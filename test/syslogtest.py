import socket
import time

def send_syslog_message_tcp(syslog_server, syslog_port, message, sdata):
    # Create a TCP socket connection to the syslog server
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    # Connect to the syslog server on the specified port
    sock.connect((syslog_server, syslog_port))
    
    # Prepare the syslog message
    # RFC 5424 format: <PRI> TIMESTAMP HOSTNAME APP-NAME PROCID MSGID [SDATA] MSG
    facility = 4  # user-level messages
    severity = 2  # critical error
    priority = (facility * 8 + severity)  # priority calculation (Facility * 8 + Severity)

    # Construct the PRI field
    pri_field = f"<{priority}>1"  # PRI is the priority field, which is <facility*8 + severity>

    # Construct timestamp (ISO 8601 format)
    timestamp = time.strftime("%Y-%m-%dT%H:%M:%S.000Z", time.gmtime())
    
    # Other syslog fields
    hostname = "myhost"
    app_name = "myapp"
    proc_id = "12345"
    msg_id = "-"  # optional, can be left as '-'
    
    msg = message  # The actual log message
    # Adding Structured Data (SDATA)
    sdata_part = f"[{sdata}]"  # SDATA part, e.g., "[exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"]"
    # Form the entire syslog message
    syslog_message = f"{pri_field} {timestamp} {hostname} {app_name} {proc_id} {msg_id} {sdata_part} {msg}\n"

    # Send the message via TCP to the syslog server
    sock.sendall(syslog_message.encode())
    
    # Close the socket connection
    sock.close()
    return syslog_message


# Example usage
syslog_server = "localhost"  # Syslog server address
syslog_port = 10312           # Syslog port (default is 514 for TCP)
log_message = "This is a test syslog message from Python"
sdata = 'exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"'  # Structured Data

print(send_syslog_message_tcp(syslog_server, syslog_port, log_message, sdata))
