# Project setup and Enabling modules

we need to let Go know that we want to use modules functionality to help manage (and version control) any third-party packages that our project imports.

But before we can do that, we need to decide what the module path for our project should be.

To do this make sure that you’re in the root of your project directory, and then run the
`go mod init command` — passing in your module path as a parameter like so:

```bash
 go mod init meeklife.net/snippetbox
```

# Web Application Basic

We will begin with three essentials:

- A handler: If you’re coming from an MVC-background, you can think of handlers as being a bit like controllers. They’re responsible for executing your application logic and for writing HTTP response headers and bodies.

- A router (servermux): This stores a mapping between the URL patterns for your application and the corresponding handlers. Usually you have one servemux for your application containing all your routes.

- A web server: One of the great things about Go is that you can establish a web server and listen for incoming requests as part of your application itself. You don’t need an external third-party server like Nginx or Apache.

Example:

```go
// Define a home handler function which writes a byte slice containing
// "Hello from meeklife" as the response body.
func homepage(w http.ResponseWriter, r *http.Request) {
     w.Write([]byte("Welcome to the hompage"))
}

// another handler
func showSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Displaying a snippet..."))
}

// another handler
func createSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Creating a new snippet"))
}

func main() {
    // Use the http.NewServeMux() function to initialize a new servemux, then
	// register the homepage function as the handler for the "/" URL pattern.
    mux := http.NewServeMux()
    mux.HandleFunc("/", homepage)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HanldeFunc("/snippet/create", createSnippet)

    // Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
    log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
```

## Fixed Path & Subtree Patterns

Go's servermux supports two different types of urls patterns; _fixed path and subtree_

Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.
With our handlers above, "/" (or another eg. "/static/" ) is a subtree path and "/snippet" is a fixed path

## Restricting the root url pattern

so if we don't want the "/" pattern to act like a catch-all. We can include a simple check in the home hander which ultimately has the same effect:

```go
func homepage(w http.ResponseWriter, r *http.Request) {
    // Check if the current request URL path exactly matches "/". If it doesn't, use
    // the http.NotFound() function to send a 404 response to the client.
    // Importantly, we then return from the handler. If we don't return the handler
    // would keep executing and also write the "Hello from SnippetBox" message.
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("Welcome to the homepage"))
}
```

## The DefaultServeMux

The `http.Handle()` and `http.HandleFunc()` functions allow us to register a route without declaring a servermux
Example:

```go
http.HandleFunc("/", home)
http.HandleFunc("/snippet", showSnippet)
```

Although it's slightly shorter, it's not recommended to use in production applications, for the sake of security, it’s generally a good idea to avoid DefaultServeMux and the corresponding helper functions. Use your own locally-scoped servemux instead, like we have been doing in this project so far.

## About RESTful Routing

It’s important to acknowledge that the routing functionality provided by Go’s servemux is pretty lightweight. It doesn’t support routing based on the request method, it doesn’t support semantic URLs with variables in them, and it doesn’t support regexp-based patterns. If you have a background in using frameworks like Rails, Django or Laravel you might find this a bit restrictive.

The reality is that Go’s servemux can still get you quite far, and for many applications is perfectly sufficient. For the times that you need more, there’s a huge choice of third-party routers that you can use instead of Go’s servemux.

# Customizing HTTP Headers

We will be updating our headers so that our route only responds to HTTP requests which use the GET, POST etc methods

We want our createsnippet() handler to respond to only POST requests so any other request method made to createsnippet() should send a 405(Method not allowed) response.

```go
func createSnippet() {
    // Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST".
    if r.Method != http.MethodPost {
        w.WriteHeader(405)
        w.Write([]byte("Method not allowed"))
        return
    }

    w.Write([]byte("Creating a new snippet"))
}
```

- It’s only possible to call w.WriteHeader() once per response, and after the status code has been written it can’t be changed. If you try to call w.WriteHeader() a second time Go will log a warning message.
- If you don’t call w.WriteHeader() explicitly, then the first call to w.Write() will automatically send a 200 OK status code to th euser.So,if you want to send a non-200 status code, you must call w.WriteHeader() before any call to w.Write().

Another improvement we can make is to include the an `Allow: POST` header with every `405 Method Not Allowed` response to let user know which method to use

```go
// use the Header().Set() to add Allow: POST
w.Header().Set("Allow", http.MethodPost)
```

## The http.Error shortcut

If you want to send a non-200 status code and a plain-text response body then it’s a good opportunity to use the http.Error() shortcut. This is a lightweight helper function which takes a given message and status code, then calls the w.WriteHeader() and w.Write() methods behind-the-scenes for us.

```go
func createSnippet() {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        // Use the http.Error() function to send a 405 status code and "Method Not
        // Allowed" string as the response body.
        http.Error(w, "Method not allowed", 405)
        return
    }
    w.Write([]byte("Creating a new snippet"))
}
```

## Manipulating the Header Map

We used w.Header().Set() to add a new header to the response header map. But there’s also Add(), Del() and Get() methods that you can use to read and manipulate the header map too.

```go
// Set a new cache-control header. If an existing "Cache-Control" header exists
// it will be overwritten.
w.Header().Set("Cache-Control", "public, max-age=31536000")

// In contrast, the Add() method appends a new "Cache-Control" header and can
// be called multiple times.
w.Header().Add("Cache-Control", "public")
w.Header().Add("Cache-Control", "max-age=31536000")

// Delete all values for the "Cache-Control" header.
w.Header().Del("Cache-Control")

// Retrieve the first value for the "Cache-Control" header.
w.Header().Get("Cache-Control")
```

# URL Query String

some url need sto accepts an id query string parameter from the user like so:
`/snippet?id`

- The handler function needs to retrieve the value of the id parameter from the URL query string, which we can do using the `r.URL.Query().Get()` method. This will always return a string value for a parameter, or the empty string "" if no matching parameter exists.

- Because the id parameter is untrusted user input, we can validate it to make sure it’s sane and sensible. For the purpose of our Snippetbox application, we want to check that it contains a positive integer value. We can do this by trying to convert the string value to an integer with the `strconv.Atoi()` function, and then checking the value is greater than zero.

Example:

```go
// Extract the value of the id parameter from the query string and try to
// convert it to an integer using the strconv.Atoi() function. If it can't
// be converted to an integer, or the value is less than 1, we return a 404 page
// not found response.
id, err := strconv.Atoi(r.URL.Query().Get("id"))
if err != nil || id < 1 {
	http.NotFound(w, r)
	return
}
// Use the fmt.Fprintf() function to interpolate the id value with our response
// and write it to the http.ResponseWriter.
fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
```

# Prroject Structure and Organization

Implement an outline structure which follows a popular and tried-and- tested approach. It’s a solid starting point, and we should be able to reuse the general structure in a wide variety of projects.

In the root of your project repository and run the following commands:

```bash
$ rm main.go
$ mkdir -p cmd/web pkg ui/html ui/static
$ touch cmd/web/main.go
$ touch cmd/web/handlers.go
```

- The cmd directory will contain the application-specific code for the executable applications in the project. For now we’ll have just one executable application — the web application — which will live under the cmd/web directory.

- The pkg directory will contain the ancillary non-application-specific code used in the project. We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project.

- The ui directory will contain the user-interface assets used by the web application. Specifically, the ui/html directory will contain HTML templates, and the ui/static directory will contain static files (like CSS and images).

There are two big benefits for using this structure:

1. It gives a clean separation between Go and non-Go assets. All the Go code we write will live exclusively under the cmd and pkg directories, leaving the project root free to hold non-Go assets like UI files, makefiles and module definitions (including our go.mod file). This can make things easier to manage when it comes to building and deploying your application in the future.

2. It scales really nicely if you want to add another executable application to your project. For example, you might want to add a CLI (Command Line Interface) to automate some administrative tasks in the future. With this structure, you could create this CLI application under cmd/cli and it will be able to import and reuse all the code you’ve written under the pkg directory.

So now our web application consists of multiple Go source code files under the cmd/web directory. To run these, we can use the go run command like so:
`go run ./cmd/web`

# HTML Templating and Inheritance

We will inject a bit of life into the project and develop a proper home page for our Snippetbox web application. Over the next couple of chapters we’ll work towards creating a page.
To do this, we will first create a new template file in the ui/html directory
`touch ui/hmtl/home.page.tmpl`
Then write our markup in the file for our homepage

So now that we’ve created a template file with the HTML markup for the home page, the next question is how do we get our home handler to render it?
For this we need to import Go’s html/template package, which provides a family of functions for safely parsing and rendering HTML templates. We can use the functions in this package to parse the template file and then execute the template.

Therefore in the handlers.go we will add these to the homepage handler;

```go
// Use the template.ParseFiles() function to read the template file into a
// template set. If there's an error, we log the detailed error message and use
// the http.Error() function to send a generic 500 Internal Server Error
// response to the user.
ts, err := template.ParseFiles("ui/html/home.page.tmpl")
if err != nil {
	log.Println(err.Error())
	http.Error(w, "Invalid Server Error", 500)
	return
}
// We then use the Execute() method on the template set to write the template
// content as the response body. The last parameter to Execute() represents any
// dynamic data that we want to pass in, which for now we'll leave as nil.
err = ts.Execute(w, nil)
if err != nil {
	log.Println(err.Error())
	http.Error(w, "Invalid Server Error", 500)
}
```

We can add partials to the template, which we need to render it as well.
So we need to update the code in the home handler to parse both templates.

```go
// Initialize a slice containing the paths to the two files. Note that the
// home.page.tmpl file must be the *first* file in the slice.
files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
}
// Use the template.ParseFiles() function to read the files and store the
// templates in a template set. Notice that we can pass the slice of file paths
// as a variadic parameter?
ts, err := template.ParseFiles(files...)
```
