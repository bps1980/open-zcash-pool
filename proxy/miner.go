package proxy

func (s *ProxyServer) processShare(cs *Session, id string, params []string) (bool, *ErrorReply) {
	// workerId := params[0]
	// jobId := params[1]
	nTime := params[2]
	extraNonce2 := params[3]
	solution := params[4]
	nOnce := cs.extraNonce1 + extraNonce2

	// submitTime := time.Unix() / 1000

	// TODO:
	// job := jobs[jobId]
	// if !job then return {Code: 21, Message: "Job not found"}

	if len(nTime) != 8 {
		return false, &ErrorReply{Code: 20, Message: "Incorrect size of nTime"}
	}

	if len(nOnce) != 64 {
		return false, &ErrorReply{Code: 20, Message: "Incorrect size of nOnce"}
	}

	if len(solution) != 2694 {
		return false, &ErrorReply{Code: 20, Message: "Incorrect size of solution"}
	}

	t := s.currentWork()

	// shareExists, validShare, errorReply := s.processShare(cs, id, t, params)
	// ok := s.policy.ApplySharePolicy(cs.ip, !shareExists && validShare)

	// if !validShare {
	// 	log.Printf("Invalid share from %s@%s", cs.login, cs.ip)
	// 	// Bad shares limit reached, return error and close
	// 	if !ok {
	// 		return false, false, errorReply
	// 	}
	// 	return false, false, nil
	// }
	// log.Printf("Valid share from %s@%s", cs.login, cs.ip)

	// if shareExists {
	// 	log.Printf("Duplicate share from %s@%s %v", cs.login, cs.ip, params)
	// 	return false, false, &ErrorReply{Code: 22, Message: "Duplicate share"}
	// }

	// if !ok {
	// 	return false, true, &ErrorReply{Code: -1, Message: "High rate of invalid shares"}
	// }
	// return false, true, nil
	return true, nil
}

// func (s *ProxyServer) processShare(cs *Session, id string, t *Work, params []string) (bool, bool, error) {
// 	nonceHex := params[0]
// 	hashNoNonce := params[1]
// 	mixDigest := params[2]
// 	nonce, _ := strconv.ParseUint(strings.Replace(nonceHex, "0x", "", -1), 16, 64)
// 	shareDiff := s.config.Proxy.Difficulty

// 	h, ok := t.headers[hashNoNonce]
// 	if !ok {
// 		log.Printf("Stale share from %v@%v", cs.login, cs.ip)
// 		return false, false, false
// 	}

// 	share := Block{
// 		number:      h.height,
// 		hashNoNonce: common.HexToHash(hashNoNonce),
// 		difficulty:  big.NewInt(shareDiff),
// 		nonce:       nonce,
// 		mixDigest:   common.HexToHash(mixDigest),
// 	}

// 	block := Block{
// 		number:      h.height,
// 		hashNoNonce: common.HexToHash(hashNoNonce),
// 		difficulty:  h.diff,
// 		nonce:       nonce,
// 		mixDigest:   common.HexToHash(mixDigest),
// 	}

// 	if !hasher.Verify(share) {
// 		return false, false, false
// 	}

// 	if hasher.Verify(block) {
// 		ok, err := s.rpc().SubmitBlock(params)
// 		if err != nil {
// 			log.Printf("Block submission failure at height %v for %v: %v", h.height, t.Header, err)
// 		} else if !ok {
// 			log.Printf("Block rejected at height %v for %v", h.height, t.Header)
// 			return false, false, false
// 		} else {
// 			s.fetchWork()
// 			exist, err := s.backend.WriteBlock(cs.login, id, params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
// 			if exist {
// 				return false, true, false
// 			}
// 			if err != nil {
// 				log.Println("Failed to insert block candidate into backend:", err)
// 			} else {
// 				log.Printf("Inserted block %v to backend", h.height)
// 			}
// 			log.Printf("Block found by miner %v@%v at height %d", cs.login, cs.ip, h.height)
// 		}
// 	} else {
// 		exist, err := s.backend.WriteShare(cs.login, id, params, shareDiff, h.height, s.hashrateExpiration)
// 		if exist {
// 			return false, true, false
// 		}
// 		if err != nil {
// 			log.Println("Failed to insert share data into backend:", err)
// 		}
// 	}
// 	return false, false, true
// }
