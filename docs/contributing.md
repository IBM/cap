# Contributing to IBM/cap

## DCO Signing
Each commit in a pull requests must be signed with a valid DCO.

First configure git for dco signing:
```
$ git config --local user.name "Full name here"
$ git config --local user.email "email address here"
```

Check your git config:
```
$ git config -l
user.name=Mike Brown
user.email=brownwm@us.ibm.com
```

Create a branch for your PR, make your changes, create your commit,
sign the commit with a signature line using the `-s` option,
and push the commit to your new branch:
```
$ git commit -s -m "this commit fixes..."
```


If you pushed your commit already, amend it with your signature then force push the amended commit:
```
$ git commit -s --amend
```

## Make Options

$ make verify
- Execute the source code verification tools"

$ make install.tools
- Install tools used by verify"
