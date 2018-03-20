tempura
=======

Deep-fried extensions for [github.com/hajimehoshi/ebiten](https://github.com/hajimehoshi/ebiten)

This is a collection of utilities for quickly creating games for Android, iOS, Javascript, and native desktop using 
[ebiten](https://github.com/hajimehoshi/ebiten). Included is a set of resources for common operations.

For real usage, check out my [Tanks](https://github.com/explodes/tanks) game for Android, Web, and Linux.

Loader
------
Loading resources is made easy with the `Loader` struct. For an in-memory cache of 
resources, you can use `CachedLoader`


Objects
-------
Most objects in games are sprites that you draw at some position with some size. The objects update every frame, 
performing actions such as movement, shooting, jumping, what have you.

`Object` encapsulates this behavior, and allows you to easily add behaviors.


Objects are all drawn an updated, so to facilitate that, `Objects` provides a way to work with groups of Objects.

Groups of objects are often draw in different layers. `Layers` makes this easy to do.


Text
----

Drawing text based on its size is a really common operation. The `Text` struct makes it easy to work with single-line
and multi-line text.
