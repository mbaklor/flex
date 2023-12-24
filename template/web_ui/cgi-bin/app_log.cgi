<% . /usr/local/lib/cgi/haserl.sh ; check_instances uinetwork.cgi %>
<% . /usr/local/lib/cgi/generic.sh ; print_http_hdr %>
<% . /usr/local/lib/cgi/status.sh %>
<% . /usr/local/lib/cgi/config.sh %>

<% 
if [ `jq 'has("app_log")' /mnt/data/package/manifest.json` = true ]; then
    APP_LOG=`jq -j .app_log /mnt/data/package/manifest.json`
else
    APP_LOG="app_log.log"
fi
%>

<html lang="en_US">
<head>
    <title>Logs</title>
    <meta http-equiv="Content-type" content="text/html;charset=UTF-8">
    <!-- <meta http-equiv="refresh" content="1"> -->
    <link href="/theme/css/main.css" rel="stylesheet" type="text/css">
</head>
<body>
    <h2>
        <% echo ${APP_LOG} %>
    </h2>
    <div class="logs-canvas" id="logs-canvas">
        <pre class="logs-canvas"><% tail -n 50 "/var/log/${APP_LOG}" %></pre>
    </div>
</body>
</html>

