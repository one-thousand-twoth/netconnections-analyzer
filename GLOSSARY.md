# Network Connection Diagrams

The language used to describe the network relationships represented by this tool.

## Language

**Declared connection**:
A source-to-destination relationship asserted by an input record. It is not a verified or observed network connection.
_Avoid_: Discovered connection, observed traffic

**Host**:
A named endpoint of a declared connection, appearing as its source or destination. The name does not by itself prove that a machine exists.
_Avoid_: Graph node, verified machine

**Host metadata**:
Information describing a host rather than any particular declared connection.
_Avoid_: Edge metadata

**Connection metadata**:
Information describing one declared connection rather than either host; ports and protocol are connection metadata.
_Avoid_: Host metadata
