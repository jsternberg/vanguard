Pipelining
==========
Server provisioning is commonly very slow because of the number of
things that must be done. One of the issues with provisioning is that
each action is performed one by one and rarely in parallel. This is
because it's hard to determine what actions might effect future actions.

Modules commonly run in two steps for idempotence. The first step is to
determine if a change must be made. The second step is to make that
change if it has to be made. We always want to err on the side of doing
something as we want to converge as quickly as possible, but there are
certain actions that are low cost and could just be done multiple times.

Let's take templating as an example. No matter what, when we issue a
template, we're going to need to template the file anyway. We can save
time by just templating the file anyway. Once we've templated the file,
it's trivial to then write that file. If we're unsure whether or not a
file needs to be changed, we should just write the file.

This comes with a lot of caveats. There are many things to consider such
as:

* What handlers/notifications will execute if I do this?
* What if this should be skipped due to something that happened earlier
  in the run?

We need to keep track of these dependencies. If a handler may get
invoked from changing that templated file, we need to wait to see if the
file really needed to be changed or not so we don't accidentally trigger
a server restart. But, if the handler would not trigger an unsafe
action, maybe we can just invoke it too without a performance hit.

The biggest point that comes out of this is modules need to communicate
more with the runtime to understand which actions can be done in advance
and which ones need to be delayed until the time the action is run.
Vanguard will attempt to make that information available to modules.
