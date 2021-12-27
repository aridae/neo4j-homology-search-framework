package KMersgenerator

type KMersGenerator struct {
	k        uint64
	OutKMers chan string
}

func NewKMersGenerator(k uint64, queueSize uint64) *KMersGenerator {
	return &KMersGenerator{
		k:        k,
		OutKMers: make(chan string, queueSize),
	}
}

// k <= 31
func (gen *KMersGenerator) Generate() error {
	defer close(gen.OutKMers)

	var x, y uint64
	var one uint64 = uint64(1)
	var kmer []byte = make([]byte, gen.k)
	alphabet := []byte{'A', 'C', 'G', 'T'}

	for x = 0; x < one<<(2*gen.k); x++ {
		y = x
		for i := uint64(0); i < gen.k; y >>= 2 {
			kmer[i] = alphabet[y&3]
			i++
		}
		gen.OutKMers <- string(kmer)
	}
	return nil
}
