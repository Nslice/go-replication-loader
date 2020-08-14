# Replication Loader

This program serves to install updates after making a backup and restarting CMS and NetPipe services.
To restart services you must have administrative privileges.

The algorithm of installation process:
* Console Monolithic Service stoppes to end of eleed sessions as is the NetPipe service to stop sync process with child databases.
* After stopping services loader starts the backup database process. Installation will not proceed without successfully created backup.