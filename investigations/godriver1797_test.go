package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

// The issue with these tests is that an undefined attribute in the JSON
// document causes the parsing to abort and fails to properly complete parsing

type response struct {
	Status  int    `bson:"status"`
	Message string `bson:"message"`
}

func createModel(doc string) (*response, error) {
	d := response{}
	err := bson.UnmarshalExtJSON([]byte(doc), false, &d)
	return &d, err
}

func Test1(t *testing.T) {
	s1 := `{
"message": "The operation was successful",
"status": 200,
"records": [{"personalIds": [2]}]
}`

	r1, err := createModel(s1)
	require.NoError(t, err)
	require.Equal(t, 200, r1.Status)
	require.Equal(t, "The operation was successful", r1.Message)
}

func Test2_Fails(t *testing.T) {
	// ONLY has issues when an undefined attribute is an array
	s2 := `{
"records": [{"personalIds": [2]}],
"message": "The operation was successful",
"status": 200
}`

	r2, err := createModel(s2)
	require.NoError(t, err)
	require.Equal(t, 200, r2.Status)
	require.Equal(t, "The operation was successful", r2.Message)
}

func Test2a(t *testing.T) {
	s2 := `{
"records": {"personalIds": 2},
"message": "The operation was successful",
"status": 200
}`

	r2, err := createModel(s2)
	require.NoError(t, err)
	require.Equal(t, 200, r2.Status)
	require.Equal(t, "The operation was successful", r2.Message)
}

func Test2b(t *testing.T) {
	s2 := `{
"records": 2,
"message": "The operation was successful",
"status": 200
}`

	r2, err := createModel(s2)
	require.NoError(t, err)
	require.Equal(t, 200, r2.Status)
	require.Equal(t, "The operation was successful", r2.Message)
}

func Test2c_Fails(t *testing.T) {
	// ONLY has issues when an undefined attribute is an array
	s2 := `{
"records": {"personalIds": [2]},
"message": "The operation was successful",
"status": 200
}`

	r2, err := createModel(s2)
	require.NoError(t, err)
	require.Equal(t, 200, r2.Status)
	require.Equal(t, "The operation was successful", r2.Message)
}

func Test3_Fails(t *testing.T) {
	// ONLY has issues when an undefined attribute is an array
	s3 := `{
"status": 200,
"records": [{"personalIds": [2]}],
"message": "The operation was successful"
}`

	r3, err := createModel(s3)
	require.NoError(t, err)
	require.Equal(t, 200, r3.Status)
	require.Equal(t, "The operation was successful", r3.Message)
}
