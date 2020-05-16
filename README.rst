PEM File Certificate Annotations
================================

Many PEM files contain several certificates without annotations. This program can
add those annotations to help identify what a given certificate block in the file 
contains.

Annotations are added above each certificate block, and look like this:

.. code-block:: text

    Subject:    CN=localhost,O=Python Software Foundation,L=Castle Anthrax,C=XY
    Issuer:     CN=localhost,O=Python Software Foundation,L=Castle Anthrax,C=XY
    Not Before: 2018-08-29 14:23:15 +00:00
    Not After:  2028-08-26 14:23:15 +00:00
    -----BEGIN CERTIFICATE-----
    MIIEWTCCAsGgAwIBAgIJAJinz4jHSjLtMA0GCSqGSIb3DQEBCwUAMF8xCzAJBgNV
    BAYTAlhZMRcwFQYDVQQHDA5DYXN0bGUgQW50aHJheDEjMCEGA1UECgwaUHl0aG9u
    ...


The program read from files specified on the command line, and outputs to standard
output. If the ``-i`` option is used, input files are edited in place, adding 
annotations. The original file will be backed up in the the same path with a ``~``
appended.

