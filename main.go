package main

import (
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/physics"
	"github.com/g3n/engine/experimental/physics/object"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	//"github.com/g3n/g3nd/app"
	"time"
)

type PhysicsSpheres struct {
	sim *physics.Simulation
	app app.Application

	sphereGeom *geometry.Geometry
	matSphere  *material.Standard

	anim   *texture.Animator
	sprite *graphic.Sprite

	attractorOn bool
	gravity     *physics.ConstantForceField
	attractor   *physics.AttractorForceField
	scene       *core.Node
}

func main() {

	// Create application and scene
	a := app.App()

	scene := core.NewNode()
	var t PhysicsSpheres
	t.sim = physics.NewSimulation(scene)
	t.gravity = physics.NewConstantForceField(&math32.Vector3{0, -0.98, 0})
	t.attractor = physics.NewAttractorForceField(&math32.Vector3{0, 1, 0}, 1)
	t.sim.AddForceField(t.gravity)
	t.scene = scene
	a.Subscribe(window.OnKeyRepeat, t.onKey)
	a.Subscribe(window.OnKeyDown, t.onKey)

	//a.Camera().GetCamera().SetPosition
	// LookAt

	// Create axes helper
	axes := helper.NewAxes(1)
	scene.Add(axes)
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

	// Create simulation and force fields

	// Decode model in in OBJ format
	dec, err := obj.Decode("gopher.obj", "gopher.mtl")
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

	sco, err := obj.Decode("one.obj", "one.mtl")
	if err != nil {
		panic(err.Error())
	}

	// Create a new node with all the objects in the decoded file and adds it to the scene
	groupO, err := sco.NewGroup()
	if err != nil {
		panic(err.Error())
	}
	groupO.SetScale(0.3, 0.3, 0.3)
	groupO.SetPosition(2, 2, -2)
	scene.Add(groupO)
	for i := range sco.Objects {
		m := sco.Objects[i]
		o, _ := sco.NewMesh(&m)
		l := object.NewBody(o)
		t.sim.AddBody(l, "Blender")
	}
	// Create and add ambient light to scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// Create and add directional white light to the scene
	dirLight := light.NewDirectional(&math32.Color{1, 1, 1}, 1.0)
	dirLight.SetPosition(1, 0, 0)
	scene.Add(dirLight)

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)
	scene.Add(t)
	t.sim.SetPaused(false)
	// Run the application
	//	app.Create().Run()
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)

	})
}

func (t *PhysicsSpheres) ThrowBall() {

	// Obtain throw direction from camera position and target
	//camTarget := t.app.Orbit().Target()

	// Create sphere rigid body
	sphere := graphic.NewMesh(t.sphereGeom, t.matSphere)
	sphere.SetPosition(1, 1, 1)
	//t.app.Scene().Add(sphere)
	t.scene.Add(sphere)
	rb := object.NewBody(sphere)
	rb.SetShape(geometry.NewSphere(.1, 1, 1))
	t.sim.AddBody(rb, "Sphere")
}

func (t *PhysicsSpheres) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	switch kev.Key {
	case window.KeyP:
		t.sim.SetPaused(!t.sim.Paused())
	case window.KeyO:
		t.sim.SetPaused(false)
		t.sim.Step(0.016)
		t.sim.SetPaused(true)
	case window.KeySpace:
		t.ThrowBall()
	case window.KeyA:
		if t.attractorOn {
			t.sim.AddForceField(t.gravity)
			t.sim.RemoveForceField(t.attractor)
			t.sprite.SetVisible(false)
			t.attractorOn = false
		} else {
			t.sim.RemoveForceField(t.gravity)
			t.sim.AddForceField(t.attractor)
			t.sprite.SetVisible(true)
			t.attractorOn = true
		}
	case window.Key2:
		// TODO
	}
}

// Update is called every frame.
func (t *PhysicsSpheres) Update(a *app.Application, deltaTime time.Duration) {

	t.sim.Step(float32(deltaTime.Seconds()))
	t.anim.Update(time.Now())
}

// Cleanup is called once at the end of the demo.
func (t *PhysicsSpheres) Cleanup(a *app.Application) {}
