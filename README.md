libvendor
=========

Libvendor allows Go applications to become vendored. Now other Go apps can use your app without worrying about your dependencies!

Due to the use of os.Getwd(), to use libvendor, you need to build the application and run the binary (go run will not suffice). 

### Usage

After ```go build libvendor.go```, all you need to do is run ```./libvendor``` in your application's root directory. 

Currently, libvendor only generates a report and tells you what files to modify and how to modify them. The next version will do all the modifications from you.
