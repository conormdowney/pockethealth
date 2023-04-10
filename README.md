# General Info

There is no front end component to the project, I just used Postman to test.

Clone the project and run main.go <br>
The server will start up and listen on localhost:3000

# What to expect

There are two endpoints:

1.) localhost:3000/upload?tag=<tagNumber>....

2.) localhost:3000/converttopng

# upload

The endpoint can accept 1 or more files and also 1 or more tags as query parameters.

The endpoint will accept multiple tags from a request like: 

localhost:3000/upload?tag=0400,0565&tag=0010,0020

Add the files to the body with the key "files"

Uploaded files are stored at C:\temp\uploads 

After a file has been sent to the service it will do the following:

1.) Save the file to C:\temp\uploads <br>
2.) Extract data for any tags that were requested. This data will be returned to the called in a map of the form

{filename: {tagNumber: {tagdata}, tagNumber2:{tagData}}, filename2: {tagNumber: {tagdata}, tagNumber2:{tagData}}

example: <br>
passing in two files, 2.dcm and 3.dcm, to the endpoint localhost:3000/upload?tag=0400,0565&tag=0010,0020 gives

{"2.dcm":{"0010,0020":{"tag":{"Group":16,"Element":32},"VR":0,"rawVR":"LO","valueLength":4,"value":["5184"]},"0400,0565":{"tag":{"Group":1024,"Element":1381},"VR":0,"rawVR":"CS","valueLength":8,"value":["CORRECT"]}},"3.dcm":{"0010,0020":{"tag":{"Group":16,"Element":32},"VR":0,"rawVR":"LO","valueLength":4,"value":["5184"]},"0400,0565":{"tag":{"Group":1024,"Element":1381},"VR":0,"rawVR":"CS","valueLength":8,"value":["CORRECT"]}}}

if there is no tag data on a given file no error will be returned. The pacakge itself returns an error. Instead it would just be left out of the map of tag data for that file.

# converttopng

The endpoint can accept 1 or more files.

Add the files to the body with the key "files"

PNG's are stored at C:\temp\pngs

After a file has been sent to the service it will do the following:

1.) Save the file to C:\temp\uploads <br>
2.) Convert the image to a png and return the name of the png file on the filesystem mapped to the original uploaded file

{filename: [png_filename], filename2: [png_filename]}

Example: 

passing 5 files to the endpoint, the result is:

{"1.dcm":["C:/temp/pngs/2c90d180-1131-40ff-aa5f-85b365c2357d_0.png"],"2.dcm":["C:/temp/pngs/599fe596-9521-4b25-8c46-8c2971c89072_0.png"],"3.dcm":["C:/temp/pngs/7e5e84df-5f02-43b1-b30d-74205cb2a384_0.png"],"4.dcm":["C:/temp/pngs/9dbe6aaf-9aaa-43df-b932-c7496fe7f1f2_0.png"],"5.dcm":["C:/temp/pngs/a9c5cde3-ff06-43c8-993e-bfa496ef27a6_0.png"]}

Files stored on the filesystem would need to be in a place that is accessible form the internet through a url, or better yet, in something like S3 in AWS or the Azure equivalent. I stored them in a folder not accessible but since there is no front end component this is really just demonstrative of what would be returned.

The .dcm file name maps to an array of strings. The images come from the parse dcms frames property. In the samples I used there was only a single frame but it is an array so I assume multiple could exist. Sending back an arrya allows for multiple images frames for a single file.

# Errors

The api will return a 400 error in the following cases:

1.) if the request body is not set as multi part form <br>
2.) If there are no files <br>
3.) If any file being uploaded is not a dcm file <br>

Any other issue that arises will return a 500

# Testing

Unit tests exist for most classes. The service class ran into some weird race issue. I think previous go routines that were created are affecting the test result becuase the issue only arises when running the test as part of the whole test suite. Running them individually works fine. 

# Code 

handler.go contains the handlers for the two endpoints. 

The validation for requests is performed here. Ensuring the issues from the Errors section are not wrong. Each handler then calls functions from the service layer.

service.go contains the business logic for the service. Upload creates a go routine for each file being uploaded that saves the file and extracts the data to return. ConverToPng saves the file and converts the image to png, returning the filenames of the pngs mapped to the associated uploaded file.

file.go handles the actual saving of the file that has been uploaded.

dicom.go is a wrapper for the dicom package. It handles calls that are required to the dicom package.

router.go is a wrapper around the package i used for handling web requests - gocraft/web https://github.com/gocraft/web
