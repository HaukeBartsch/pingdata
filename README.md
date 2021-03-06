Download DICOM images from XNAT
===============================

Command-line program connects to an XNAT image store, lists the subjects available and downloads subject
information as zip files to the local directory. This program requires a valid login on the image store.

Example:

	> pingdata list
	SubjID, Age, Gender
	# subjects found: 782
	...

	> pingdata pull P0007
	download all DICOM files for subject P0007
	[283.55mb] P0007.zip

You can download the compiled program for your platform here. On Linux and MacOS make the downloaded file executable with 'chmod +x pingdata'.

* MacOSX:
	https://raw.github.com/HaukeBartsch/pingdata/master/binary/MacOSX/pingdata
* Linux64:
	https://raw.github.com/HaukeBartsch/pingdata/master/binary/Linux64/pingdata
* Windows:
	https://raw.github.com/HaukeBartsch/pingdata/master/binary/Windows/pingdata.exe

Help page (pingdata help):

	NAME:
	  pingdata - Download PING data from http://www.nitrc.org/ir.
	
	This program uses the XNAT REST API to download subject data. Start by listing
	the subjects available:
	
	> pingdata list
	 
	To download the data of a specific subject call:
	
	> pingdata pull PXXXXX
	
	where PXXXXX is the subject identification number.
	
	USAGE:
	  pingdata [global options] command [command options] [arguments...]
	
	VERSION:
	  0.0.1
	  
	AUTHOR:
	  Hauke Bartsch - <HaukeBartsch@gmail.com>
	  
	COMMANDS:
	  pull, p	Retrieve subject data as zip
	  list, l	Retrieve list of subjects
	  help, h	Shows a list of commands or help for one command
	  
	GLOBAL OPTIONS:
	  --help, -h		show help
	  --version, -v	print the version
