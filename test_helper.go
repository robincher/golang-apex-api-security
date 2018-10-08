package apexsigner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

// GOLANG ()
const GOLANG string = "golang"

// DEFAULT ()
const DEFAULT string = "default"

// TestParam ()
type TestParam struct {
	ID          string `json:"id"`
	Description string `json:"description"`

	APIParam APIParam `json:"apiParam"`

	PublicCertFileName string      `json:"publicCertFileName"`
	SkipTest           []string    `json:"skipTest"`
	Message            string      `json:"message"`
	Debug              bool        `json:"debug"`
	TestTag            bool        `json:"testTag"`
	ExpectedResult     interface{} `json:"expectedResult"`
	ErrorTest          bool        `json:"errorTest"`
}

type testParam interface {
	getID() string
	getDescription() string
	getErrorTest() bool
	getExpectedResult() string
	setExpectedResult(string)
}

type fx func(*TestParam) (string, error)

func executeTest(t *testing.T, fileName string, testFunction fx) {
	testParamsArray, err := getTestParamFromJSONArray(fileName)
	if err != nil {
		t.Errorf("\nparameters file error::: %s", err.Error())
	}

	testCaseCount := 0

	for _, param := range testParamsArray {
		if !contains(param.SkipTest, GOLANG) {
			testCaseCount++

			result, err := testFunction(&param)

			handleMessage(&param, result, err, t)
		}
	}

	if testCaseCount == 0 {
		t.Errorf("\nno test executed...")
	}

	t.Logf(">>> %d test executed <<<", testCaseCount)
}

func (testParam *TestParam) getID() string {
	return testParam.ID
}

func (testParam *TestParam) getDescription() string {
	return testParam.Description
}

func (testParam *TestParam) getErrorTest() bool {
	return testParam.ErrorTest
}

func (testParam *TestParam) getExpectedResult() string {

	switch expectedResult := testParam.ExpectedResult.(type) {
	case string:
		return expectedResult
	case map[string]interface{}:
		if g, ok := expectedResult[GOLANG].(string); ok {
			return g
		} else if d, ok := expectedResult[DEFAULT].(string); ok {
			return d
		}
	}
	return ""
}

func (testParam *TestParam) setExpectedResult(newValue string) {
	switch expectedResult := testParam.ExpectedResult.(type) {
	case map[string]interface{}:
		if expectedResult[GOLANG] != nil {
			expectedResult[GOLANG] = newValue
		} else if expectedResult[DEFAULT] != nil {
			expectedResult[DEFAULT] = newValue
		}
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func getInvokeURL(param APIParam) (string, error) {
	url, err := url.Parse(param.InvokeURL)
	if err != nil {
		return "", err
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return "", fmt.Errorf("not supported URL scheme(%s), only http and https allow", url.Scheme)
	}

	query := url.Query()

	formDataArray := parseParam(param.QueryString)

	for _, values := range formDataArray {
		query.Add(values[0], values[1])
	}

	url.RawQuery = query.Encode()

	return url.String(), nil
}

func handleMessage(param testParam, result string, err error, t *testing.T) {
	if param.getErrorTest() {
		if err == nil {
			t.Error(formatMessage(param, param.getExpectedResult(), "but no error raise"))
		} else if err.Error() != param.getExpectedResult() {
			t.Error(formatMessage(param, param.getExpectedResult(), err.Error()))
		}
	} else if err != nil {
		t.Error(formatMessage(param, "no error raise", err.Error()))
	} else if result != param.getExpectedResult() {
		t.Error(formatMessage(param, param.getExpectedResult(), result))
	}
}

func formatMessage(param testParam, expect string, actual string) string {
	var message string
	message = message + fmt.Sprintf("\nID     ::: %s", param.getID())
	message = message + fmt.Sprintf("\ndesc   ::: %s", param.getDescription())
	message = message + fmt.Sprintf("\nexpect ::: %s", expect)
	message = message + fmt.Sprintf("\nbut got::: %s", actual)

	return message
}

func getTestParamFromJSONArray(fileName string) ([]TestParam, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(fileName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Printf("Successfully Opened %s\n", fileName)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	// we initialize our Users array
	var params []TestParam

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &params)

	return params, nil
}
