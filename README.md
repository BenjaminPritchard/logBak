#About
Example program for working with Google embedded NoSQL engine LevelDB from a GoLang program. 

#Description
This simple utility can convert back and forth between a plain text log file, and a LevelDB database file. 

#Usage
	logBak b log_file NoSQL_file		-> backups log_file to NoSQL_file
	logBak r log_file NoSQL_file		-> restores log_file from NoSQL_file
