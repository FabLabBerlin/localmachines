#User client

Please use `io.js`. At the time of writing this the latest version of io.js is `3.0.0`.

##Quickstart

- `npm test` to run unit tests
- `npm run-script dev` to build the dev code on demand "watch" (`dev/bundle.js` and other files in `prod/`)
- `npm run-script prod` to build the production code (`prod/bundle.js` and other files in `prod/`)

##Convention used:

###General

- The achitecture of the src files are :
  - src/
    - js/
     - JsFiles.js
    - assets/
      - css/scss/less
      - img/..
- All function are in **lowerCamelCase**
- All component **name** and **const** are in **UpperCamelCase**
- Test directories are "__test__"
- Test file are : *NameOfYourFile*-**test**.js

###Router

We use react-router.  Refer to https://github.com/rackt/react-router/blob/0.13.x/docs/guides/flux.md for details on how to use it in our Flux architecture.

###Components

A component is a class with a `render()` function.  The class has a state
and the render function is only called when the state changes.

Here an idea how the code is organize:

 - Common parts
   - [mixins](https://facebook.github.io/react/docs/component-specs.html#mixins)
   - [static](https://facebook.github.io/react/docs/component-specs.html#statics)
   - [`getInitialState()`](https://facebook.github.io/react/docs/component-specs.html#getinitialstate)
   - *stuff*
   - `onChange`
   - `render()`: required method
 - Lifecycle methods (needed for low-level DOM access)
   - [`componentWillMount()`](https://facebook.github.io/react/docs/component-specs.html#mounting-componentwillmount): called before DOM is built
   - [`componentDidMount()`](https://facebook.github.io/react/docs/component-specs.html#mounting-componentdidmount): called after DOM is built
   - [`componentWillUnmount()`](https://facebook.github.io/react/docs/component-specs.html#unmounting-componentwillunmount)

The function in *stuff* aren't organize in a special way, up to you to specify your convention here

###Stores

There are 2 stores in this interface:
 - LoginStore
 - MachineStore

The **Loginstore** manages the **login phase**, and is used essentialy by Login and LoginNfc component

The **MachineStore** manages the data the main page needs to be display correctly.

Some convention are done for the store because it's a really complicate file. It needs to be respected for making the editing easier:
 - On top, in commentary, put all the main functions
 - This is how the file should be organized: 
   - state
   - api call
   - function related to api call (for formating the state for example)
   - getter
   - onChange
 - This organisation will be kept until a real flux is implemented
 - Function name for **api call** are to be : **apiMETHODFunctionName**
  - For example: *apiGetActivationActive*
 - **getter** name has to start with get and follow by the Information you get(for example: *getIsLogged*).

###Flux

We're trying to use a flux architecture even if there is no dispatcher.
The application there is too small to bother using one.

For this there is some rules:
 - Component can't access to any store directly, they have to call action
 - Make as stateless component as possible
 - Try to regroup the state in some major component to debug easily
 - All the interaction with the back-end is done in the store
 - We try to make it a circle: **actions =>** *dispatcher* **=> store => view => action**

##Architecture

```
  LoginStore                MachineStore
      |     \________________      |
      |                      \     |
 LoginChooser                 MachinePage
      |                            |
      |                       MachineList
     / \                           |
Login   LoginNfc              MachineChooser
                                   |
                                   |
                                  /|\
                                 / | \
                                /  |  \
                         Occupied Busy Free
                            \    /
                             \  /
                             Timer

```
                        
###Directories

- Components are in component folder
- Components are calling actions
- Actions are in actions folder
- Actions are calling store methods
- Stores are in store folder
- Store changes components state(at least Top component state) to update the view
- Test folder are in each folder you want to test (if you want to test store: src/js/store/__test__)

##Thank

Thank you for reading it and try your best to write clean code and stick to the convention (you can change it to fit the team better).
