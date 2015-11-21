# logsclient

Logsclient is a tool for downloading log files from CloudFlare's Enterprise Log Share. 

## Usage

-auth.email="": authorization email

-auth.key="": authorization key

-dir="/var/folders/j5/y4pm_7yj1qd1_7dh2kxdf5fm0000gn/T/": directory to download logs into

-end=1448147588: the unix epoch timestamp to end downloading at

-max=72h0m0s: the maximum time in the past the start can be

-start=-1: the unix epoch timestamp to start downloading from

-url="https://api.cloudflare.com/client/v4/logs": URL for CloudFlare logs API
