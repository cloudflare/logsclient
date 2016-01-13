# logsclient

Logsclient is a tool for downloading log files from CloudFlare's Enterprise Log Share. 

## Usage

  -auth.email string
        authorization email
  -auth.key string
        authorization key
  -dir string
        directory to download logs 
  -end int
        the unix epoch timestamp to end downloading at (defaults to time the program is run )
  -interval duration
        the time interval to save files in (default 1m0s)
  -max duration
        the maximum time in the past the start can be (default 72h0m0s)
  -start int
        the unix epoch timestamp to start downloading from (default -1 and looks for the checkpoint file in the download directory)
  -url string
        URL for CloudFlare logs API - https://api.cloudflare.com/client/v4/zones/<zone tag>/logs/requests
