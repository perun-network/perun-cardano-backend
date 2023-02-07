package types

import "perun.network/go-perun/channel"

const (
	DepositedTag  = "Deposited"
	DisputedTag   = "Disputed"
	ProgressedTag = "Progressed"
	ConcludedTag  = "Concluded"
	WithdrawnTag  = "Withdrawn"
)

// TODO: Figure out what to return on AdjudicatorEvent.Version(), if there is no concept of a state version for that event

type (
	Deposited struct {
		ChannelID    ID
		ChannelDatum ChannelDatum
	}
	Disputed struct {
		ChannelID    ID
		ChannelDatum ChannelDatum
	}
	Progressed struct {
		ChannelID    ID
		ChannelDatum ChannelDatum
	}
	Concluded struct {
		ChannelID ID
	}
	Withdrawn struct {
		ChannelID ID
	}
)

func (w Withdrawn) ID() channel.ID {
	return w.ChannelID
}

func (w Withdrawn) Timeout() channel.Timeout {
	return nil
}

func (w Withdrawn) Version() uint64 {
	return 0
}

func (c Concluded) ID() channel.ID {
	return c.ChannelID
}

func (c Concluded) Timeout() channel.Timeout {
	return nil
}

func (c Concluded) Version() uint64 {
	return 0
}

func (p Progressed) ID() channel.ID {
	return p.ChannelID
}

func (p Progressed) Timeout() channel.Timeout {
	return &channel.TimeTimeout{Time: p.ChannelDatum.Time.Add(p.ChannelDatum.ChannelParameters.Timeout)}
}

func (p Progressed) Version() uint64 {
	return p.ChannelDatum.ChannelState.Version
}

func (d Disputed) ID() channel.ID {
	return d.ChannelID
}

func (d Disputed) Timeout() channel.Timeout {
	return &channel.TimeTimeout{Time: d.ChannelDatum.Time.Add(d.ChannelDatum.ChannelParameters.Timeout)}
}

func (d Disputed) Version() uint64 {
	return d.ChannelDatum.ChannelState.Version
}

func (d Deposited) ID() channel.ID {
	return d.ChannelID
}

func (d Deposited) Timeout() channel.Timeout {
	return nil
}

func (d Deposited) Version() uint64 {
	return 0
}
