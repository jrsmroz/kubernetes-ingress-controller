//go:build conformance_tests
// +build conformance_tests

package conformance

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	"sigs.k8s.io/gateway-api/conformance/tests"
	"sigs.k8s.io/gateway-api/conformance/utils/suite"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/controllers/gateway"
)

var (
	showDebug     = true
	shouldCleanup = true

	manifestRepo                  = "https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/master/"
	conformanceTestsBaseManifests = fmt.Sprintf("%s/conformance/base/manifests.yaml", manifestRepo)
)

func TestGatewayConformance(t *testing.T) {
	t.Parallel()

	t.Log("configuring environment for gateway conformance tests")
	client, err := client.New(env.Cluster().Config(), client.Options{})
	require.NoError(t, err)
	require.NoError(t, gatewayv1alpha2.AddToScheme(client.Scheme()))

	t.Log("creating GatewayClass for gateway conformance tests")
	gwc := &gatewayv1alpha2.GatewayClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid.NewString(),
		},
		Spec: gatewayv1alpha2.GatewayClassSpec{
			ControllerName: gateway.ControllerName,
		},
	}
	require.NoError(t, client.Create(ctx, gwc))
	defer func() {
		require.NoError(t, client.Delete(ctx, gwc))
	}()

	t.Log("starting the gateway conformance test suite")
	cSuite := suite.New(suite.Options{
		Client:               client,
		GatewayClassName:     gwc.Name,
		Debug:                showDebug,
		CleanupBaseResources: shouldCleanup,
		BaseManifests:        conformanceTestsBaseManifests,
		SupportedFeatures: []suite.SupportedFeature{
			suite.SupportReferenceGrant,
		},
	})
	cSuite.Setup(t)

	t.Log("configuring gateway conformance tests")
	for i := range tests.ConformanceTests {
		for j, manifest := range tests.ConformanceTests[i].Manifests {
			tests.ConformanceTests[i].Manifests[j] = fmt.Sprintf("%s/conformance/%s", manifestRepo, manifest)
		}
	}

	t.Log("running gateway conformance tests")
	for _, tt := range tests.ConformanceTests {
		if enabledGatewayConformanceTests.Has(tt.ShortName) {
			t.Run(tt.Description, func(t *testing.T) { tt.Run(t, cSuite) })
		}
	}
}

var enabledGatewayConformanceTests = sets.NewString(
	"GatewaySecretInvalidReferenceGrant",
	"GatewaySecretMissingReferenceGrant",
	"GatewaySecretReferenceGrantAllInNamespace",
	"GatewaySecretReferenceGrantSpecific",
)
