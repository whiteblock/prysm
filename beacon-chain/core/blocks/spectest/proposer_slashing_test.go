package spectest

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
)

func runProposerSlashingTest(t *testing.T, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not load file %v", err)
	}

	test := &ProposerSlashingMinimal{}
	if err := yaml.Unmarshal(file, test); err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}

	if err := spectest.SetConfig(test.Config); err != nil {
		t.Fatal(err)
	}

	for _, tt := range test.TestCases {
		t.Run(tt.Description, func(t *testing.T) {
			pre := &pb.BeaconState{}
			if err := convertToPb(tt.Pre, pre); err != nil {
				t.Fatal(err)
			}

			expectedPost := &pb.BeaconState{}
			if err = convertToPb(tt.Post, expectedPost); err != nil {
				t.Fatal(err)
			}

			proposerSlashing := &pb.ProposerSlashing{}
			if err = convertToPb(tt.ProposerSlashing, proposerSlashing); err != nil {
				t.Fatal(err)
			}

			block := &pb.BeaconBlock{Body: &pb.BeaconBlockBody{ProposerSlashings: []*pb.ProposerSlashing{proposerSlashing}}}

			var postState *pb.BeaconState
			postState, err = blocks.ProcessProposerSlashings(pre, block, true)
			// Note: This doesn't test anything worthwhile. It essentially tests
			// that *any* error has occurred, not any specific error.
			if len(expectedPost.ValidatorRegistry) == 0 {
				if err == nil {
					t.Fatal("Did not fail when expected")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(postState, expectedPost) {
				t.Error("Post state does not match expected")
			}
		})
	}
}

func TestProposerSlashingMinimal(t *testing.T) {
	runProposerSlashingTest(t, "proposer_slashing_minimal.yaml")
}

func TestProposerSlashingMainnet(t *testing.T) {
	runProposerSlashingTest(t, "proposer_slashing_mainnet.yaml")
}
