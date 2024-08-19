package main

type AmCommand struct {
	Command string `json:"command"`
}

type AmPool struct {
	Pool                int    `json:"POOL"`
	URL                 string `json:"URL"`
	Status              string `json:"Status"`
	User                string `json:"User"`
	CurrentBlockVersion int    `json:"Current Block Version"`
}

type AmReplySummary struct {
	Status []struct {
		STATUS      string `json:"STATUS"`
		When        int    `json:"When"`
		Code        int    `json:"Code"`
		Msg         string `json:"Msg"`
		Description string `json:"Description"`
	} `json:"STATUS"`
	Summary []struct {
		Elapsed            int     `json:"Elapsed"`
		GHS5S              float64 `json:"GHS 5s"`
		GHSAv              float64 `json:"GHS av"`
		GHS30M             float64 `json:"GHS 30m"`
		FoundBlocks        int     `json:"Found Blocks"`
		Getwork            int     `json:"Getwork"`
		Accepted           int     `json:"Accepted"`
		Rejected           int     `json:"Rejected"`
		HardwareErrors     int     `json:"Hardware Errors"`
		Utility            float64 `json:"Utility"`
		Discarded          int     `json:"Discarded"`
		Stale              int     `json:"Stale"`
		GetFailures        int     `json:"Get Failures"`
		LocalWork          int     `json:"Local Work"`
		RemoteFailures     int     `json:"Remote Failures"`
		NetworkBlocks      int     `json:"Network Blocks"`
		TotalMH            int64   `json:"Total MH"`
		WorkUtility        float64 `json:"Work Utility"`
		DifficultyAccepted float64 `json:"Difficulty Accepted"`
		DifficultyRejected float64 `json:"Difficulty Rejected"`
		DifficultyStale    float64 `json:"Difficulty Stale"`
		BestShare          int64   `json:"Best Share"`
		DeviceHardware     float64 `json:"Device Hardware%"`
		DeviceRejected     float64 `json:"Device Rejected%"`
		PoolRejected       float64 `json:"Pool Rejected%"`
		PoolStale          float64 `json:"Pool Stale%"`
		LastGetwork        int     `json:"Last getwork"`
	} `json:"SUMMARY"`
}

type AmReplyStats struct {
	Status []CgMinerStatus
	Stats  []struct {
		BMMiner        string  `json:"BMMiner,omitempty"`
		Miner          string  `json:"Miner,omitempty"`
		CompileTime    string  `json:"CompileTime,omitempty"`
		Type           string  `json:"Type,omitempty"`
		STATS          int     `json:"STATS,omitempty"`
		ID             string  `json:"ID,omitempty"`
		Elapsed        int     `json:"Elapsed,omitempty"`
		Calls          int     `json:"Calls,omitempty"`
		Wait           int     `json:"Wait,omitempty"`
		Max            int     `json:"Max,omitempty"`
		Min            int     `json:"Min,omitempty"`
		GHS5S          float64 `json:"GHS 5s,omitempty"`
		GHSAv          float64 `json:"GHS av,omitempty"`
		Rate30M        float64 `json:"rate_30m,omitempty"`
		Mode           int     `json:"Mode,omitempty"`
		MinerCount     int     `json:"miner_count,omitempty"`
		Frequency      int     `json:"frequency,omitempty"`
		FanNum         int     `json:"fan_num,omitempty"`
		Fan1           int     `json:"fan1,omitempty"`
		Fan2           int     `json:"fan2,omitempty"`
		Fan3           int     `json:"fan3,omitempty"`
		Fan4           int     `json:"fan4,omitempty"`
		TempNum        int     `json:"temp_num,omitempty"`
		Temp1          int     `json:"temp1,omitempty"`
		Temp21         int     `json:"temp2_1,omitempty"`
		Temp2          int     `json:"temp2,omitempty"`
		Temp22         int     `json:"temp2_2,omitempty"`
		Temp3          int     `json:"temp3,omitempty"`
		Temp23         int     `json:"temp2_3,omitempty"`
		TempPcb1       string  `json:"temp_pcb1,omitempty"`
		TempPcb2       string  `json:"temp_pcb2,omitempty"`
		TempPcb3       string  `json:"temp_pcb3,omitempty"`
		TempPcb4       string  `json:"temp_pcb4,omitempty"`
		TempChip1      string  `json:"temp_chip1,omitempty"`
		TempChip2      string  `json:"temp_chip2,omitempty"`
		TempChip3      string  `json:"temp_chip3,omitempty"`
		TempChip4      string  `json:"temp_chip4,omitempty"`
		TempPic1       string  `json:"temp_pic1,omitempty"`
		TempPic2       string  `json:"temp_pic2,omitempty"`
		TempPic3       string  `json:"temp_pic3,omitempty"`
		TempPic4       string  `json:"temp_pic4,omitempty"`
		TotalRateideal int     `json:"total_rateideal,omitempty"`
		RateUnit       string  `json:"rate_unit,omitempty"`
		TotalFreqavg   int     `json:"total_freqavg,omitempty"`
		TotalAcn       int     `json:"total_acn,omitempty"`
		TotalRate      float64 `json:"total rate,omitempty"`
		TempMax        int     `json:"temp_max,omitempty"`
		NoMatchingWork int     `json:"no_matching_work,omitempty"`
		ChainAcn1      int     `json:"chain_acn1,omitempty"`
		ChainAcn2      int     `json:"chain_acn2,omitempty"`
		ChainAcn3      int     `json:"chain_acn3,omitempty"`
		ChainAcn4      int     `json:"chain_acn4,omitempty"`
		ChainAcs1      string  `json:"chain_acs1,omitempty"`
		ChainAcs2      string  `json:"chain_acs2,omitempty"`
		ChainAcs3      string  `json:"chain_acs3,omitempty"`
		ChainAcs4      string  `json:"chain_acs4,omitempty"`
		ChainHw1       int     `json:"chain_hw1,omitempty"`
		ChainHw2       int     `json:"chain_hw2,omitempty"`
		ChainHw3       int     `json:"chain_hw3,omitempty"`
		ChainHw4       int     `json:"chain_hw4,omitempty"`
		ChainRate1     string  `json:"chain_rate1,omitempty"`
		ChainRate2     string  `json:"chain_rate2,omitempty"`
		ChainRate3     string  `json:"chain_rate3,omitempty"`
		ChainRate4     string  `json:"chain_rate4,omitempty"`
		Freq1          int     `json:"freq1,omitempty"`
		Freq2          int     `json:"freq2,omitempty"`
		Freq3          int     `json:"freq3,omitempty"`
		Freq4          int     `json:"freq4,omitempty"`
		MinerVersion   string  `json:"miner_version,omitempty"`
		MinerID        string  `json:"miner_id,omitempty"`
	} `json:"STATS"`
}

type AmReplyPools struct {
	Status []CgMinerStatus `json:"STATUS"`
	Pools  []AmPool        `json:"POOLS"`
}
