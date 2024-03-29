"""
Q: How many authenticated connections from the Automation Agent are accepted and ended in the mongos logs included in HELP-29818?
A: 295
"""
import json
import datetime

logpath = "./help29818_logs/atlas-b3e2ep-shard-00-00.azg38.mongodb.net/mongos/27016/mongodb/mongodb.log"

conn_states = {}
conn_accepted = 0
conn_ended = 0
automation_agent_conn_authed = 0
automation_agent_conn_authed_ended = 0
automation_agent_conn_authed_accepted_and_ended = 0
automation_agent_conn_authed_accepted_and_ended_durations = []

def init_state (connid):
    conn_states[connid] = {
        "state": "",
        "authed": False,
        "from_automation_agent": False,
        "time_accepted": None,
        "history": []
    }

def check_state (connid, *expected_states):
    global conn_states
    if connid not in conn_states:
        # Always initialize if connid is not in conn_states.
        # Connections may be in any state before the log capture.
        init_state (connid)
        return

    if conn_states[connid]["state"] not in expected_states:
        raise Exception ("Expected connection {} to be in one of states: {} but got {}".format(connid, expected_states, conn_states[connid]["state"]))

with open (logpath, "r") as logfile:
    for rawline in logfile:
        if rawline.strip() == "":
            continue

        line = json.loads(rawline)
        if line["msg"] == "***** SERVER RESTARTED *****":
            print ("See log line '{}'. Resetting connection states".format(line["msg"]))
            conn_states = {}
        if line["msg"] == "Connection accepted":
            connid = "conn" + str(line["attr"]["connectionId"])
            if connid in conn_states:
                raise Exception ("Expected connection {} not to be in conn_states, but got {}".format(connid, conn_states[connid]))
            init_state (connid)
            conn_states[connid]["history"].append (rawline)
            conn_states[connid]["state"] = "Connection accepted"
            conn_states[connid]["accepted"] = True
            conn_states[connid]["time_accepted"] = datetime.datetime.fromisoformat(line["t"]["$date"])
            conn_accepted += 1
            continue

        if line["msg"] == "client metadata":
            connid = line["ctx"]
            check_state (connid, "Connection accepted")
            conn_states[connid]["history"].append (rawline)
            conn_states[connid]["state"] = "client metadata"
            if line["attr"]["doc"]["application"]["name"].startswith("MongoDB Automation Agent v11.8.1.7231"):
                conn_states[connid]["from_automation_agent"] = True
            continue

        if line["msg"] == "Authentication succeeded":
            connid = line["ctx"]
            check_state (connid, "client metadata")
            conn_states[connid]["history"].append (rawline)
            conn_states[connid]["state"] = "Authentication succeeded"
            conn_states[connid]["authed"] = True
            if conn_states[connid]["from_automation_agent"]:
                automation_agent_conn_authed += 1
            continue

        if line["msg"] == "Connection ended":
            connid = line["ctx"]
            check_state (connid, "Connection accepted", "client metadata", "Authentication succeeded")
            conn_states[connid]["history"].append (rawline)
            conn_ended += 1
            if conn_states[connid]["authed"] and conn_states[connid]["from_automation_agent"]:
                automation_agent_conn_authed_ended += 1
                if conn_states[connid]["accepted"]:
                    automation_agent_conn_authed_accepted_and_ended += 1
                    time_ended = datetime.datetime.fromisoformat (line["t"]["$date"])
                    duration = time_ended - conn_states[connid]["time_accepted"]
                    if connid[4] != '4' and connid[4] != '5' and duration.total_seconds() < 1.5:
                        print ("check out connection: {}".format(connid))
                        print ("".join(conn_states[connid]["history"]))
                        import sys
                        sys.exit(1)
                    automation_agent_conn_authed_accepted_and_ended_durations.append(duration.total_seconds())
            del conn_states[connid]
            continue

print ("conn_accepted={}".format(conn_accepted))
print ("conn_ended={}".format(conn_ended))
print ("automation_agent_conn_authed={}".format(automation_agent_conn_authed))
print ("automation_agent_conn_authed_ended={}".format(automation_agent_conn_authed_ended))
print ("automation_agent_conn_authed_accepted_and_ended={}".format(automation_agent_conn_authed_accepted_and_ended))
avg = sum(automation_agent_conn_authed_accepted_and_ended_durations) / len (automation_agent_conn_authed_accepted_and_ended_durations)
automation_agent_conn_authed_accepted_and_ended_durations.sort()
median = automation_agent_conn_authed_accepted_and_ended_durations[len(automation_agent_conn_authed_accepted_and_ended_durations) // 2]
print ("automation_agent_conn_authed_accepted_and_ended_durations average: {}".format(avg))
print ("automation_agent_conn_authed_accepted_and_ended_durations median: {}".format(median))

"""
Output:
conn_accepted=2271
conn_ended=2223
automation_agent_conn_authed=354
automation_agent_conn_authed_ended=295
automation_agent_conn_authed_accepted_and_ended=295
automation_agent_conn_authed_accepted_and_ended_durations average: 62.94652881355929
automation_agent_conn_authed_accepted_and_ended_durations median: 1.05

Interpretation:
First log line includes this timestamp: 2021-12-24T17:07:04.834+00:00
Last log line includes this timestamp:  2021-12-24T17:58:20.486+00:00
Time span is about 50 minutes, or 3000 seconds.
Majority of application connections from Automation Agent are short lived (~1 second).
Observed ~300 application connections created / closed from Automation Agent.
Suggests that Automation Agent is creating short-lived mongo.Clients.
"""