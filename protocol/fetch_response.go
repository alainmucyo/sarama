package protocol

type FetchResponseBlock struct {
	Err                 KError
	HighWaterMarkOffset int64
	MsgSet              MessageSet
}

func (pr *FetchResponseBlock) decode(pd packetDecoder) (err error) {
	pr.Err, err = pd.getError()
	if err != nil {
		return err
	}

	pr.HighWaterMarkOffset, err = pd.getInt64()
	if err != nil {
		return err
	}

	msgSetSize, err := pd.getInt32()
	if err != nil {
		return err
	}

	msgSetDecoder, err := pd.getSubset(int(msgSetSize))
	if err != nil {
		return err
	}
	err = (&pr.MsgSet).decode(msgSetDecoder)

	return err
}

type FetchResponse struct {
	Blocks map[string]map[int32]*FetchResponseBlock
}

func (fr *FetchResponse) decode(pd packetDecoder) (err error) {
	numTopics, err := pd.getArrayCount()
	if err != nil {
		return err
	}

	fr.Blocks = make(map[string]map[int32]*FetchResponseBlock, numTopics)
	for i := 0; i < numTopics; i++ {
		name, err := pd.getString()
		if err != nil {
			return err
		}

		numBlocks, err := pd.getArrayCount()
		if err != nil {
			return err
		}

		fr.Blocks[name] = make(map[int32]*FetchResponseBlock, numBlocks)

		for j := 0; j < numBlocks; j++ {
			id, err := pd.getInt32()
			if err != nil {
				return err
			}

			block := new(FetchResponseBlock)
			err = block.decode(pd)
			if err != nil {
				return err
			}
			fr.Blocks[name][id] = block
		}
	}

	return nil
}
