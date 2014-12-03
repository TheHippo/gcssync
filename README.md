# gcssync

**TOC:**
<!-- MarkdownTOC autolink=true bracket=round -->

- [Usage](#usage)

<!-- /MarkdownTOC -->


Syncs files with a Google Cloud Storage bucket.

When syncing it checks for modifcation and update time to prevent uploading files in the next run.


## Usage

    NAME:
       gcssync - Sync files with Google Cloud Storage
    
    USAGE:
       gcssync [global options] command [command options] [arguments...]
    
    VERSION:
       0.0.0
    
    COMMANDS:
       list, l  List remote files
       upload, u    Upload a single file
       sync, s  Syncs a folder to a Google Cloudstorage bucket
       help, h  Shows a list of commands or help for one command
       
    GLOBAL OPTIONS:
       --cachefile 'cache.json' Cache file for caching auth tokens [$AUTH_CACHE_FILE]
       --bucketname, -b         Name of bucket [$BUCKET_NAME]
       --projectid, -p      Google project [$PROJECT_ID]
       --clientid, -c       Auth client id [$AUTH_CLIENT_ID]
       --clientsecret, -s       Client secrect [$AUTH_CLIENT_SECRET]
       --code           Authorization Code [$AUTH_CODE]
       --help, -h           show help
       --version, -v        print the version

