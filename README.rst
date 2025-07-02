PEM File Certificate Annotations
================================

Certificates in DER format, while compatible with the ASCII character are
opaque to human readers. While ``openssl x509 -text -in somefile.crt``
provides human-readable information, it is sub-optimal for some use-cases:

*  Only handles one certificate at a time.
*  Information overload!
*  Ugly formatting.

This tool can handle multiple files and multiple certificates per file,
and outputs the most important information just before each certificate block:

.. code-block:: text

    $ annotate-pem reg.crt | head
    Subject:          CN=kreg,OU=Research division,O=Mouse Inc.,L=Floating in space
    Issuer:           CN=INT CA,OU=Research Division,O=Mouse Inc.,L=Floating in space
    Not Before:       2025-07-01 14:58:44 +00:00
    Not After:        2027-07-11 14:58:44 +00:00
    Subj. Alt. Names: DNS:kreg, DNS:kreg.local, IP:192.168.122.215
    -----BEGIN CERTIFICATE-----
    MIIIBDCCBCygAwIBAgICEAwwDQYJKoZIhvcNAQELBQAwXjEaMBgGA1UEBwwRRmxv
    YXRpbmcgaW4gc3BhY2UxEzARBgNVBAoMCk1vdXNlIEluYy4xGjAYBgNVBAsMEVJl
    c2VhcmNoIERpdmlzaW9uMQ8wDQYDVQQDDAZJTlQgQ0EwHhcNMjUwNzAxMTQ1ODQ0
    WhcNMjcwNzExMTQ1ODQ0WjBcMRowGAYDVQQHDBFGbG9hdGluZyBpbiBzcGFjZTET


The program reads from files specified on the command line, and outputs to
standard output. If the ``-i`` option is used, input files are edited in place,
adding annotations (the original file will be backed-up with ``~`` appended).

