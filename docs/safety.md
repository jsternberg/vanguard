Safety
======
One of the most important functions of cluster management is determining
safety of a particular action. Most updates can be done in the
background and do not require any downtime of a service. Some updates
need a service restart and require the server to incur downtime. Updates
of high availability services either need to be done all at once on all
servers (requiring downtime) or can be done on each server individually
and not require any downtime for the cluster as a whole.

Vanguard attempts to solve this problem by giving each task a set of
needs and promises. Here are a few examples:

Installing a new package
------------------------
Imagine we have a server and it was functioning previously. We now want
a new shell, [fish](https://github.com/fish/fish), to be installed on
every server. This package will likely not cause any functional problem
with the server itself if it were installed. It's not going to restart
apache or require a change to any configuration files that are being
used by the active service. It's merely an additional package that gets
installed.

Most installs are like this. They aren't going to cause any downtime for
an existing part of the system and can be done in the background with no
repurcussions. This would be an example of a *safe* action and could be
done during a non-downtime period.

Upgrading a critical package
----------------------------
Not all package upgrades are created equal. For this one, we are talking
about package upgrades for things like `libssl`, `apache2`, or some
other critical service used by your application.

These upgrades would incur downtime for that specific server and need to
be delayed until a suitable maintenance window. We need to issue a
*promise* that we will upgrade the package. We may also have other
tasks, such as a service restart, that want to be done after the
upgrade. For these, they issue a *need* for the service to be upgraded
and get delayed.
