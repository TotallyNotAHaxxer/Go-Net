```

                                   
 _____ _____     _____ _____ _____ 
|   __|     |___|   | |   __|_   _|
|  |  |  |  |___| | | |   __| | |  
|_____|_____|   |_|___|_____| |_|  
                                   

```

# example usage 

First create a directory called `Modules` then create a file outside of modules called main.go 

add the source code files in this repo into modules and write the following code. For example the code below uses a Modules dir and impliments the proxy 

```go
package main

import Network "main/Modules"

func main() {
	Network.Proxy()
}

```
