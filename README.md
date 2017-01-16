# quiet-dog

A silent doorbell that emails and texts when an Amazon Dash button is pressed

## Configuration

Create a file "config.json" in the same directory as the executable with the following shape:

```javascript
{
    "mac": "The MAC address of your dash button",
    "twilioKey": "The Twilio API key for sending texts",
    "phoneNumbers": ['phone numbers', 'to text when people are here'],
    "interface": "The network interface to listen for ARP packets on (probably wlan0 or similar)"
}
```

## Running

Run `go build` then `sudo ./quiet-dog`