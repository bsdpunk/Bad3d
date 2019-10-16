package main

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/physics"
	"github.com/g3n/engine/experimental/physics/object"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
	"os"
	"time"
)

type solid struct {
	Name string
	sim  *physics.Simulation
	anim *texture.Animator
}

func main() {

	// Create application and scene
	a := app.App()
	scene := core.NewNode()

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// Decode model in in OBJ format
	dec, err := obj.Decode(os.Args[1]+".obj", os.Args[1]+".mtl")
	if err != nil {
		panic(err.Error())
	}

	// Create a new node with all the objects in the decoded file and adds it to the scene
	group, err := dec.NewGroup()
	if err != nil {
		panic(err.Error())
	}
	group.SetScale(0.3, 0.3, 0.3)
	group.SetPosition(0.0, -0.8, -0.2)
	scene.Add(group)
	sim := physics.NewSimulation(scene)

	gravity := physics.NewConstantForceField(&math32.Vector3{0, -0.98, 0})
	//attractor := physics.NewAttractorForceField(&math32.Vector3{0, 1, 0}, 1)
	sim.AddForceField(gravity)
	//sim.AddForceField(attractor)

	for i := range dec.Objects {
		m := dec.Objects[i]
		o, _ := dec.NewMesh(&m)
		l := object.NewBody(o)
		l.SetVelocity(&math32.Vector3{0, 0, float32(i)})
		fmt.Printf("%+v\n", l)
		sim.AddBody(l, "Blender"+string(i))
	}
	fmt.Printf("%+v\n", sim)
	fmt.Printf("%+v\n", scene)
	fmt.Printf("%+v\n", sim.Bodies())
	// Create and add ambient light to scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// Create and add directional white light to the scene
	dirLight := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	dirLight.SetPosition(1, 0, 0)
	scene.Add(dirLight)

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
	var phy solid
	phy.Name = "orig"
	phy.sim = sim
	tex, err := texture.NewTexture2DFromImage("smoke30.png")
	phy.anim = texture.NewAnimator(tex, 6, 5)
	phy.anim.SetDispTime(2 * 16666 * time.Microsecond)
	// Run the application
	//scene.Add(phy.anim)
	//phy.Update(a, time.Duration(1000))
	//:w
	//phy.Update(a, time.Duration)
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
		//		phy.Update(a, time.Now())
		//phy.Update(a, deltaTime)
	})
}

func (phy *solid) Update(a *app.Application, deltaTime time.Duration) {

	phy.sim.Step(float32(deltaTime.Seconds()))
	phy.anim.Update(time.Now())
	//phy.Update(a, ,time.Now())
	fmt.Println("Hey")
}
