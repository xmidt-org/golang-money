package money

// SubTracerContainer holds the necessary objects for a
// subtracer process.
type SubTracerContainer struct {
	function  SubTracer
	htChannel chan<- *HTTPTracker
}

// StarterContainer holds the necessary objects for a
// starter process.
type StarterContainer struct {
	function  Starter
	htChannel chan<- *HTTPTracker
}

// Ender holds the necessary objects for a
// ender process.

type EnderContainer struct {
	function Ender
}
