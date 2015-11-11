package txtedit

var EVERYTHING_CONF string = `
# ====== ntpd ======
restrict 127.0.0.1
discard average 3 minimum 1

# ====== httpd ======
DirectoryIndex index.html index.html.var
<Files ~ "^\.ht">
    <IfModule mod_access_compat.c>
        Order allow,deny
        Deny from all
    </IfModule>
</Files>

# ======= named =======
zone "." in {
    type hint;
    file "root.hint";
    forwarders { 192.0.2.1; 192.0.2.2; };
};

# ===== postfix =====
debugger_command =
     PATH=/bin:/usr/bin:/usr/local/bin:/usr/X11R6/bin
     ddd $daemon_directory/$process_name $process_id & sleep 5
relay_domains = $mydestination, hash:/etc/postfix/relay

# ===== vsftpd =====
pam_service_name=vsftpd

# ===== dhcpd =====
class "foo" {
  match if substring (option vendor-class-identifier, 0, 4) = "SUNW";
}

shared-network 224-29 {
  subnet 10.17.224.0 netmask 255.255.255.0 {
    allow members of "foo";
    range 10.17.224.10 10.17.224.250;
  }
}

# ===== inetd =====
defaults
{
log_type        = FILE /var/log/xinetd.log
}

# ===== xl2tpd =====
[lns default]
ip range = 192.168.1.128-192.168.1.254
; listen-addr = 192.168.1.98
`
