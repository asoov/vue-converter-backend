package helpers

import (
	"testing"
)

func TestExtractJavaScriptFromHTML(t *testing.T) {
	// Define a test case
	testHTML := `<html>
<head>
    <script type="text/javascript">
    console.log("Hello, world!");
    </script>
    <title>Test Page</title>
</head>
<body>
    <script>
    alert("Another script");
    </script>
</body>
</html>`

	expected := `
	console.log("Hello, world!");
            
	alert("Another script");
	
`

	// Call the function
	result := ExtractJavaScriptFromHTML(testHTML)

	// Check the result
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
