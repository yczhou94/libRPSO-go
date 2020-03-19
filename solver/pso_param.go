package solver

type PSOParam struct {
	w             float64
	c1            float64
	c2            float64
	c3            float64
	pr            float64
	pm            float64
	t             float64
	bound         *Bound
	dim           int
	nProc         int
	popSize       int
	maxStep       int
	simAnnealFlag bool
	initFunc      InitParticleFunc
	targetFunc    TargetFunc
	args          interface{}
}

func (p *PSOParam) SetW(w float64) {
	p.w = w
}

func (p *PSOParam) SetC1(c1 float64) {
	p.c1 = c1
}

func (p *PSOParam) SetC2(c2 float64) {
	p.c2 = c2
}

func (p *PSOParam) SetC3(c3 float64) {
	p.c3 = c3
}

func (p *PSOParam) SetPr(pr float64) {
	p.pr = pr
}

func (p *PSOParam) SetPm(pm float64) {
	p.pm = pm
}

func (p *PSOParam) SetT(t float64) {
	p.t = t
}

func (p *PSOParam) SetBound(bound *Bound) {
	p.bound = bound
}

func (p *PSOParam) SetDim(dim int) {
	p.dim = dim
}

func (p *PSOParam) SetNProc(nProc int) {
	p.nProc = nProc
}

func (p *PSOParam) SetInitFunc(initFunc InitParticleFunc) {
	p.initFunc = initFunc
}

func (p *PSOParam) SetTargetFunc(targetFunc TargetFunc) {
	p.targetFunc = targetFunc
}

func (p *PSOParam) SetTargetFuncArgs(args interface{}) {
	p.args = args
}

func (p *PSOParam) SetSimAnnealFlag(flag bool) {
	p.simAnnealFlag = flag
}

func NewPSOParam(popSize int, dim int, targetFunc TargetFunc) *PSOParam {
	return &PSOParam{
		w:  0.723,
		c1: 1.4454,
		c2: 1.4454,
		c3: 0.72,
		pr: 0.01,
		pm: 0.01,
		t:  100,
		bound: &Bound{
			XUpper: 10,
			XLower: -10,
			VUpper: 5,
			VLower: -5,
		},
		dim:           dim,
		nProc:         1,
		popSize:       popSize,
		simAnnealFlag: true,
		initFunc:      DefaultInitParticle,
		targetFunc:    targetFunc,
		args:          nil,
	}
}
