# Instructions

The api I made can accept 1 or more files and also 1 or more tags as query parameters. 
The api will accept multiple tags from a request like: 

localhost:3000/upload?tag=0400,0565&tag=0010,0020

There is no front end component to the project, I just used Postman to test.

Uploaded files are stored at C:\temp\uploads
PNG's are stored at C:\temp\pngs

Clone the project and run main.go
The server will start up and listen on localhost:3000


After a file has been sent to the service it will do the following:

1.) Save the file to C:\temp\uploads
2.) Extract data for any tags that were requested. This data will be returned to the called in a map of the form

{filename: {tagNumber: {tagdata}, tagNumber2:{tagData}}, filename2: {tagNumber: {tagdata}, tagNumber2:{tagData}}

example:
passing in two files, 2.dcm and 3.dcm, to the endpoint localhost:3000/upload?tag=0400,0565&tag=0010,0020 gives

{"2.dcm":{"0010,0020":{"tag":{"Group":16,"Element":32},"VR":0,"rawVR":"LO","valueLength":4,"value":["5184"]},"0400,0565":{"tag":{"Group":1024,"Element":1381},"VR":0,"rawVR":"CS","valueLength":8,"value":["CORRECT"]}},"3.dcm":{"0010,0020":{"tag":{"Group":16,"Element":32},"VR":0,"rawVR":"LO","valueLength":4,"value":["5184"]},"0400,0565":{"tag":{"Group":1024,"Element":1381},"VR":0,"rawVR":"CS","valueLength":8,"value":["CORRECT"]}}}

The api will return a 400 error in the following cases:

1.) if the request body is not set as multi part form
2.) If there are no files
3.) If any file being uploaded is not a dcm file

Unit tests exist for most classes. The service class ran into some weird race issue. I think previous go routines that were created are affecting the test result becuase the issue only arises when running the test as part of the whole test suite. Running them individually works fine. 

