"""
Q: How many authenticated connections from the Automation Agent are opened and closed in the mongos logs included in HELP-29818?
A: 
"""
import json

logpath = "./help29818_logs/atlas-b3e2ep-shard-00-00.azg38.mongodb.net/mongos/27016/mongodb/mongodb.log"

conn_states = {}
conn_accepted = 0
conn_ended = 0
automation_agent_conn_authed = 0
automation_agent_conn_authed_ended = 0

automation_agent_conn_authed_accepted_and_ended = 0

def init_state (connid):
    conn_states[connid] = {
        "state": "",
        "authed": False,
        "from_automation_agent": False,
        "history": [],
        "accepted": False
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

lineno = 1

with open (logpath, "r") as logfile:
    for rawline in logfile:
        if rawline.strip() == "":
            continue
    
        lineno += 1
        line = json.loads(rawline)
        if line["msg"] == "***** SERVER RESTARTED *****":
            print ("See log line '{}'. Resetting connection states".format(line["msg"]))
            conn_states = {}
        if line["msg"] == "Connection accepted":
            connid = "conn" + str(line["attr"]["connectionId"])
            if connid in conn_states:
                raise Exception ("Expected connection {} not to be in conn_states, but got {}".format(connid, conn_states[connid]))
            init_state (connid)
            conn_states[connid]["state"] = "Connection accepted"
            conn_states[connid]["history"].append(rawline)
            conn_states[connid]["accepted"] = True
            conn_accepted += 1
            continue

        if line["msg"] == "client metadata":
            connid = line["ctx"]
            check_state (connid, "Connection accepted")
            conn_states[connid]["history"].append(rawline)
            conn_states[connid]["state"] = "client metadata"
            if line["attr"]["doc"]["application"]["name"].startswith("MongoDB Automation Agent v11.8.1.7231"):
                conn_states[connid]["from_automation_agent"] = True
            continue

        if line["msg"] == "Authentication succeeded":
            connid = line["ctx"]
            check_state (connid, "client metadata")
            conn_states[connid]["history"].append(rawline)
            conn_states[connid]["state"] = "Authentication succeeded"
            conn_states[connid]["authed"] = True
            if conn_states[connid]["from_automation_agent"]:
                automation_agent_conn_authed += 1
            continue

        if line["msg"] == "Connection ended":
            connid = line["ctx"]
            check_state (connid, "Connection accepted", "client metadata", "Authentication succeeded")
            conn_states[connid]["history"].append(rawline)
            conn_ended += 1
            if conn_states[connid]["authed"] and conn_states[connid]["from_automation_agent"]:
                automation_agent_conn_authed_ended += 1
                if conn_states[connid]["accepted"]:
                    automation_agent_conn_authed_accepted_and_ended += 1
            del conn_states[connid]
            continue

print ("conn_accepted={}".format(conn_accepted))
print ("conn_ended={}".format(conn_ended))
print ("automation_agent_conn_authed={}".format(automation_agent_conn_authed))
print ("automation_agent_conn_authed_ended={}".format(automation_agent_conn_authed_ended))
print ("automation_agent_conn_authed_accepted_and_ended={}".format(automation_agent_conn_authed_accepted_and_ended))

"""
Output:


First log line includes this timestamp: 2021-12-24T17:07:04.834+00:00
Last log line includes this timestamp:  2021-12-24T17:58:20.486+00:00
Time span is about 50 minutes, or 3000 seconds.
Observed about 300 application connections created / closed.
"""