# пока не рабочий вариант
WITH 8 as k, "ACGTGTCCGATGACTTG" as input, 8 as position
WITH  size(input) - k + 1 as kmerscnt, k, input, position
UNWIND range(position - k, position + k - 1) AS X
WITH left(right(input, kmerscnt + k - 1 - X), k) AS tmpkmer, kmerscnt, k
MATCH (selectedNode:KMer {value:tmpkmer})
WITH collect(selectedNode) as selectedNodes, k
WITH selectedNodes, selectedNodes[0] as head, last(selectedNodes) as tail, k
CALL apoc.cypher.run("MATCH path=(h)-[:Precedes*"+(k + 1)+"]->(t) RETURN path;", {h: head, t:tail}) YIELD value as paths
UNWIND paths as path
return custom.pathToString(path); 

# найти все варианты во входной строке в 8й позиции
PROFILE 
WITH 8 as k, "AAAAAAAGTGTCGTCGTCAAAAA" as input, 8 as position
WITH  size(input) - k + 1 as kmerscnt, k, input, position
UNWIND range(position - k, position + k - 1) AS X
WITH left(right(input, kmerscnt + k - 1 - X), k) AS tmpkmer, kmerscnt, k
MATCH (selectedNode:KMer {value:tmpkmer})
WITH collect(selectedNode) as selectedNodes, k
WITH selectedNodes, selectedNodes[0] as head, last(selectedNodes) as tail, k
MATCH (s:Sequence)
MATCH (g:Genome)<-[:Belongs]-(s)
MATCH path=(head)-[:Precedes*8..9]->(tail) 
WHERE ALL(n in nodes(path) WHERE (n)-[:Belongs]->(s))
return {
    genome: g.name, 
    sequence: s.name, 
    variant: custom.pathToString(path), 
    location: {start_kmer: ID(head), end_kmer: ID(tail)}
};