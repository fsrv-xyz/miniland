# miniland

*A simple golang excluseive userland for linux systems.*

## learnings

### Applications are not able to listen on all addresses

Listening on all interfaces require having a minimal `/etc/hosts` file. Only entries for
IPv4 and IPv6 localhost are reqired to have http services listening on all addresses with `:<port>` pattern.