# JavaScript Frontend Documentation

## React

The (new) Frontend is based on [React](https://facebook.github.io/react/) which allows to modularize the code using JavaScript classes. Each class has a render() function that returns the ["XHTML"](https://en.wikipedia.org/wiki/XHTML) to be displayed. The combination is called [JSX](https://jsx.github.io/) and usually it is a [pure function](https://en.wikipedia.org/wiki/Pure_function).

## ImmutableJS

Application state is as usual in Functional Programming (mostly) [immutable](https://en.wikipedia.org/wiki/Immutable_object). [ImmutableJS](https://facebook.github.io/immutable-js/) is a JavaScript library for that. When opening the Browser console, the current (global) [state](https://en.wikipedia.org/wiki/State_(computer_science)) can be seen.

## Flux

[Flux](https://facebook.github.io/flux/) is basically uni-directional data flow combined with [MVC](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller). Upon input from the User or the Network, the global state transitions to a new state. The state changes from the UI are triggered by onClick, ... handlers. These then call the actions which in turn dispatch the actual state change to the store(s).

## NuclearJS

[NuclearJS](https://github.com/optimizely/nuclear-js) uses a completely immutable global state and getters must be used to get data from the stores. A getter can source stores/other getters which means there can be hierarchically cascaded. (["Functional lens"](https://medium.com/@dtipson/functional-lenses-d1aba9e52254)) Two levels of hierarchy could be compared to Model vs ViewModel like in [MVVM](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93viewmodel) but offers even more levels of abstraction if needed.
