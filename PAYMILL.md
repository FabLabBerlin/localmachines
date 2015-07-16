# Paymill integration notice
A sum up of all informations about paymill and how to integrate it.

# API library
First, there was no integration for paymill API in golang but I found out this repository : [eaigner/paymill](https://github.com/eaigner/paymill). So I forked it and now you can find it here : [KevinBacas/paymill](https://github.com/KevinBacas/paymill).

The documentation is on the README of this repository.
## What do we need to do on the API
- Move it to v2.1
- Documentation
- Use it properly and maybe add some fonctionnalities

# What Paymill can do ?
- Manage clients
- Manage payment methods
- Manage transactions
- Manage subscriptions
- Everything is accessible through a public [API](https://api.paymill.com/v2.1/)

# What do we need ?
- Real SSL certificate (User trust, security, reliability, warranty)
- Register real account onto Paymill
- Know how it's supposed to work when a payment method expires and a subscription is on

# What we need to know ?
- Registration time is now longer :
  - User side : Have to fill another form for the payment method
  - Server side : Have to send HTTP calls to Paymill (We can't run without internet FTM (Maybe a buffering system in case ?))
- Business part is more **reliable** (**No more** *Someone forgot to add this user on fastbill!*)

# What is done for the moment ?
- Created a branch ```feature-paymill```
- Automatic user add to paymill when signed up
- Forked ```eaigner/paymill``` into ```KevinBacas/paymill```
- Wrote down some documentation on how to use it in **Golang**
- API call for public key
- Two step form
