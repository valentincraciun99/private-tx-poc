package zk
import (
	"fmt"
	"github.com/arnaucube/go-snark"
	"github.com/arnaucube/go-snark/circuitcompiler"
	"math/big"
	"strings"
)

type Setup struct{
	vk snark.Vk
	pk snark.Pk
}

func GenerateCircuit() *circuitcompiler.Circuit {
	flatCode := `
func exp3(private a):
	b = a * a
	c = a * b
	return c

func main(private s0, public s1):
	s3 = exp3(s0)
	s4 = s3 + s0
	s5 = s4 + 5
	equals(s1, s5)
	out = 1 * 1
`
	parser := circuitcompiler.NewParser(strings.NewReader(flatCode))
	circuit, _ := parser.Parse()


	return circuit
}

func GenerateWitness(private int, public int, circuit *circuitcompiler.Circuit) []*big.Int {
	x := big.NewInt(int64(private))
	privateInputs := []*big.Int{x}
	y := big.NewInt(int64(public))
	publicSignals := []*big.Int{y}

	w, _ := circuit.CalculateWitness(privateInputs, publicSignals)

	return w
}

func GenerateSetup(witness []*big.Int, circuit *circuitcompiler.Circuit) (snark.Setup, []*big.Int) {
	fmt.Println("generating R1CS from flat code")
	a, b, c := circuit.GenerateR1CS()

	alphas, betas, gammas, _ := snark.Utils.PF.R1CSToQAP(a, b, c)

	_, _, _, px := snark.Utils.PF.CombinePolynomials(witness, alphas, betas, gammas)

	setup, _ := snark.GenerateTrustedSetup(len(witness), *circuit, alphas, betas, gammas)

	return setup,px
}

