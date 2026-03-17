# Certificates

### Install the Certificate Authority (CA) certificate

Install the `ca.crt` as a trusted root certificate authority into your keychain using the standard import function in MacOS keychain. 

### Update the trust settings for the CA certificate

Once installed you'll specifically need to change it's state to "Always trust" yourself. It should ask you for your machine password before saving this change.

### Check chrome certificate register

Check [Chorme Certificate Manager](chrome://certificate-manager/localcerts/platformcerts) to make sure the trusted CA root is present

### Use the localhost key and certificate

The `localhost.crt` and `localhost.key` files have already been signed with the above `ca.crt` and can be used in our microservices for dev purposes.

### Additional reference

https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/