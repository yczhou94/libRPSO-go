package libRPSO

type TargetFunc func(x []float64, args ...interface{}) ([]float64, float64, error)

type Solver struct {
	s    *swarm
	conf *SolverConf
}

type SolverConf struct {
	MaxStep     int
	NTerm       int
	PrintEvery  int
	Convergence float64
}

func NewSolver(psoParam *PSOParam, conf *SolverConf) (*Solver, error) {
	var err error
	solver := &Solver{
		conf: conf,
	}
	solver.s, err = newSwarm(psoParam)

	return solver, err
}

func (s *Solver) Run() error {
	for i := 0; i < s.conf.MaxStep; i++ {
		err := s.s.step()
		if err != nil {
			return err
		}
	}

	return nil
}
