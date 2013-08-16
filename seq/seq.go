// The seq container is to provide and manage basic information about a biological sequence.
package seq


// GenericSeq provides basic fields that all Seq type should contain. 
type GenericSeq struct {
	 DisplayId string
	 Accession string
	 PrimaryId string
	 Id string
	 Length int
	 Description string
	 IsCircular bool
}

// A container to hold various arguments pass to different methods
type Options struct {
	 Terminator string
	 Unknown string
	 Frame string
	 CodonTableId string
	 CompleteCodons bool
	 Throw bool
	 Complete bool
	 Orf string
	 Start string
	 Offset int
}

// An interface to provide information about sequence
type SeqInformer interface {
	 Seq() string
	 Subseq() string
	 Alphabet() string
}

// An interface to provide a transformed representation of GenericSeq
type SeqTransformer interface {
	 Revcom() (*GenericSeq, error)
	 Trunc(start int, end int) (*GenericSeq, error)
	 Translate(opt *Options) (*GenericSeq, error)

}
