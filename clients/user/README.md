#User client

##Script

Three scripts are implemented to help you develop or test the application:
 - **npm test**
  - It will launch jest and test all js file with **-test**
  - The configurations are in npm and preprocessor.js
 - **npm run-script dev**
  - It will launch the application in dev mode on :8080
  - watch is activated (--watch)
  - bundle.js is produce in dev/
 - **npm start**
  - It will launch the application in prod
  - bundle.js is produce in prod/

##Convention used:

###Components

Here an idea how the code is organize:
 - mixins
 - static
 - getInitialState
 - *stuff*
 - onChange
 - componentDidMount
 - Render

The function in *stuff* aren't organize in a special way, up to you to specify your convention here

###Flux

We're trying to use a flux architecture even if there is no dispatcher.
The application there is too small to bother using one.

For this there is some rules:
 - Component can't access to any store directly, they have to call action
 - Make as stateless component as possible
 - Try to regroup the state in some major component to debug easily
 - All the interaction with the back-end is done in the store
 - We try to make it a circle: **actions =>** *dispatcher* **=> store => view => action**

##Contribution

Fell free to add some rules/scripts and explain it there.
