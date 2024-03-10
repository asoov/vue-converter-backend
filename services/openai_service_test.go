package services

import "testing"

var getVueChatCompletionInput = "<template>\n  <div id=\"app\">\n    <h1>{{ message }}</h1>\n    <button @click=\"reverseMessage\">Reverse Message</button>\n  </div>\n</template>\n\n<script>\nexport default {\n  data() {\n    return {\n      message: 'Hello Vue!'\n    }\n  },\n  methods: {\n    reverseMessage() {\n      this.message = this.message.split('').reverse().join('')\n    }\n  }\n}\n</script>\n\n<style scoped>\n#app {\n  font-family: 'Avenir', 'Helvetica', 'Arial', sans-serif;\n  -webkit-font-smoothing: antialiased;\n  -moz-osx-font-smoothing: grayscale;\n  text-align: center;\n  color: #2c3e50;\n  margin-top: 60px;\n}\nh1 {\n  font-size: 2em;\n}\nbutton {\n  font-size: 1em;\n  padding: 0.5em 1em;\n  border: none;\n  background-color: #42b983;\n  color: white;\n  cursor: pointer;\n  border-radius: 3px;\n  outline: none;\n}\nbutton:hover {\n  background-color: #68d5a4;\n}\n</style>"

var getVueChatCompletionExpected = "This code is VueJS code that is not Version 3 and is not using the composition API. Please change this code that composition API is implemented. Just return the code, no explaining text. This is the code" + getVueChatCompletionInput + "\n" + "Always make sure the whole component is returned including template, script and style. If you come across special properties prefixed with a '$' make sure to destructure it from the context parameter in the setup function."

func TestGetVueChatCompletionMessage(t *testing.T) {
	result := GetVueChatCompletion(getVueChatCompletionInput)
	if result != getVueChatCompletionExpected {
		t.Errorf("GetVueChatCompletion() = %v, want %v", result, getVueChatCompletionExpected)
	}
}
