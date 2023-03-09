package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/edaniels/golog"
	"github.com/golang/geo/r3"
	"github.com/pkg/errors"
	commonpb "go.viam.com/api/common/v1"
	"go.viam.com/rdk/components/arm/universalrobots"
	"go.viam.com/rdk/components/arm/xarm"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

var scene *config
var sceneFS referenceframe.FrameSystem
var logger golog.Logger = golog.NewLogger("omplbindings")

const testArmFrame = "arm"
const testEndEffectorFrame = "end_effector"

type config struct {
	Start            []referenceframe.Input
	Goal             spatialmath.Pose
	RobotFrame       referenceframe.Frame
	EndEffectorFrame *referenceframe.Frame
	WorldState       *referenceframe.WorldState
}

var allScenes = map[int]func() (*config, error){
	// Basic scenes begin
	1:  scene1,
	2:  scene2,
	3:  scene3,
	4:  scene4,
	5:  scene5,
	6:  scene6,
	7:  scene7,
	8:  scene8,
	9:  scene9,
	10: scene10,
	11: scene11,
	12: scene12,
	// User scenes begin
	13: objSearch,
}

// initScene takes a scene number and loads the relevant information into memory
func initScene(sceneNum int) (err error) {
	if sceneFn, ok := allScenes[sceneNum]; ok {
		scene, err = sceneFn()
		if err != nil {
			return
		}
		sceneFS = referenceframe.NewEmptySimpleFrameSystem("test")
		sceneFS.AddFrame(scene.RobotFrame, sceneFS.World())

		if scene.EndEffectorFrame != nil {
			sceneFS.AddFrame(*(scene.EndEffectorFrame), scene.RobotFrame)
		}

		return
	}
	return errors.Errorf("scene %d does not exist", sceneNum)
}

func scene1() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 100
	goalPt.Y += 100

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       &referenceframe.WorldState{},
	}, nil
}

func scene2() (*config, error) {
	model, _ := xarm.Model("arm", 7)
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 200
	goalPt.Z += 100

	// Obstacles
	testPose := spatialmath.NewPose(
		r3.Vector{X: 1., Y: -200., Z: 3.},
		&spatialmath.R4AA{Theta: 0, RX: 0., RY: 0., RZ: 1.},
	)
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: referenceframe.World,
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(testPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 2000,
								Y: 2000,
								Z: 20,
							}},
						},
					},
					{
						Center: spatialmath.PoseToProtobuf(testPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 2000,
								Y: 20,
								Z: 2000,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

func scene3() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := r3.Vector{X: -400, Y: 350, Z: 0}

	// Obstacles
	testPose := spatialmath.NewPose(
		r3.Vector{X: 0., Y: 150., Z: 0.},
		&spatialmath.R4AA{Theta: 0, RX: 0., RY: 0., RZ: 1.},
	)
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(testPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 2000,
								Y: 20,
								Z: 120,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

func scene4() (*config, error) {
	model, _ := xarm.Model("arm", 6)
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 300

	// Obstacles
	testPt := startPose.Point()
	testPt.X += 150
	testPose := spatialmath.NewPose(testPt, startPose.Orientation())
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(testPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 20,
								Y: 2000,
								Z: 60,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

func scene5() (*config, error) {
	model, _ := xarm.Model("arm", 7)
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 400

	// Obstacles
	wallPose := spatialmath.NewPoseFromPoint(r3.Vector{0, -200, 0})
	obs1Pose := spatialmath.NewPoseFromPoint(r3.Vector{300, 0, 0})
	obs2Pose := spatialmath.NewPoseFromPoint(r3.Vector{300, 0, 500})
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(wallPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 2000,
								Y: 50,
								Z: 2000,
							}},
						},
					},
					{
						Center: spatialmath.PoseToProtobuf(obs1Pose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 50,
								Y: 1250,
								Z: 200,
							}},
						},
					},
					{
						Center: spatialmath.PoseToProtobuf(obs2Pose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 50,
								Y: 1250,
								Z: 200,
							}},
						},
					},
				},
			},
		},
	})

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, err
}

func scene6() (*config, error) {
	cfg, err := scene5()
	if err != nil {
		return nil, err
	}
	obstacle, err := spatialmath.NewBox(spatialmath.NewPoseFromPoint(r3.Vector{-150, 0, 0}), r3.Vector{20, 2000, 2000}, "")
	if err != nil {
		return nil, err
	}
	cfg.WorldState.Obstacles = append(
		cfg.WorldState.Obstacles,
		referenceframe.NewGeometriesInFrame(referenceframe.World, []spatialmath.Geometry{obstacle}),
	)
	return cfg, err
}

func scene7() (*config, error) {
	cfg, err := scene4()
	if err != nil {
		return nil, err
	}
	left_wall, err := spatialmath.NewBox(spatialmath.NewPoseFromPoint(r3.Vector{0, 140, 0}), r3.Vector{2000, 20, 2000}, "left_wall")
	if err != nil {
		return nil, err
	}
	right_wall, err := spatialmath.NewBox(spatialmath.NewPoseFromPoint(r3.Vector{0, -140, 0}), r3.Vector{2000, 20, 2000}, "right_wall")
	if err != nil {
		return nil, err
	}
	cfg.WorldState.Obstacles = append(
		cfg.WorldState.Obstacles,
		referenceframe.NewGeometriesInFrame(referenceframe.World, []spatialmath.Geometry{left_wall, right_wall}),
	)
	return cfg, nil
}

func scene8() (*config, error) {
	cfg, err := scene2()
	if err != nil {
		return nil, err
	}
	cfg.Goal = spatialmath.NewPose(cfg.Goal.Point(), &spatialmath.R4AA{Theta: 0, RX: 0., RY: 0., RZ: 1.})
	return cfg, err
}

func scene9() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{0, 0, 0, 0, 0, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 1100
	goalPt.Y += 600

	rGen := rand.New(rand.NewSource(int64(1)))
	obstacles := make([]spatialmath.Geometry, 0)
	for i := 0; i < 100; i++ {
		cubePose := spatialmath.NewPoseFromPoint(r3.Vector{
			X: 2000 * (rGen.Float64() - 0.5),
			Y: 2000 * (rGen.Float64() - 0.5),
			Z: 2000 * (rGen.Float64() - 0.5),
		})
		label := strconv.Itoa(i)
		cube, err := spatialmath.NewBox(cubePose, r3.Vector{1, 1, 1}, label)
		if err != nil {
			return nil, err
		}
		obstacles = append(obstacles, cube)
	}

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState: &referenceframe.WorldState{Obstacles: []*referenceframe.GeometriesInFrame{
			referenceframe.NewGeometriesInFrame(referenceframe.World, obstacles),
		}},
	}, nil

}

func scene10() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{0, -math.Pi / 4, math.Pi / 2, 3 * math.Pi / 4, -math.Pi / 2, 0})
	startPose, _ := model.Transform(startInput)

	// Goal pose
	goalPt := startPose.Point()
	goalPt.X += 1200
	goalPt.Y += 600

	// Pose of UR5 mount pillar
	pillarPose := spatialmath.NewPose(
		r3.Vector{X: 0., Y: 0., Z: -1000.},
		&spatialmath.R4AA{Theta: 0, RX: 1., RY: 0., RZ: 0.},
	)
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(pillarPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 130,
								Y: 130,
								Z: 2000,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             spatialmath.NewPose(goalPt, startPose.Orientation()),
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

// Corresponds to move that has been demonstrated to cause a self-collision on the UR5's basic planning
func scene11() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{3.8141, -1.3106, 2.4543, 4.9485, -3.4041, -2.6749})

	// Goal pose
	goalPos := r3.Vector{X: -244.43, Y: -255.12, Z: 676.97}
	goalRot := spatialmath.R3ToR4(r3.Vector{X: 0.233, Y: -1.637, Z: 1.224})
	goalPose := spatialmath.NewPose(goalPos, goalRot)

	// Pose of UR5 mount pillar
	pillarPose := spatialmath.NewPose(
		r3.Vector{X: 0., Y: 0., Z: -1000.},
		&spatialmath.R4AA{Theta: 0, RX: 1., RY: 0., RZ: 0.},
	)
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(pillarPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 130,
								Y: 130,
								Z: 2000,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             goalPose,
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

// Corresponds to move that only works with MoveJ from an engineering move set
func scene12() (*config, error) {
	model, _ := universalrobots.Model("arm")
	startInput := referenceframe.FloatsToInputs([]float64{1.2807, -1.4437, -1.3287, 3.7446, 1.4315, -0.2135})

	// Goal pose
	goalPos := r3.Vector{X: -50.47, Y: -366.47, Z: 189.04}
	goalRot := spatialmath.R3ToR4(r3.Vector{X: 0.808, Y: 2.168, Z: 2.916})
	goalPose := spatialmath.NewPose(goalPos, goalRot)

	// Pose of UR5 mount pillar
	pillarPose := spatialmath.NewPose(
		r3.Vector{X: 0., Y: 0., Z: -1000.},
		&spatialmath.R4AA{Theta: 0, RX: 1., RY: 0., RZ: 0.},
	)
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(pillarPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 130,
								Y: 130,
								Z: 2000,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             goalPose,
		RobotFrame:       model,
		EndEffectorFrame: nil,
		WorldState:       worldState,
	}, nil
}

// Corresponds to a user application involving unstructured visual search with an end effector
func objSearch() (*config, error) {
	model, _ := universalrobots.Model("arm")
	// Taken from recordings of recreations at the NYC-1900 robotics lab
	startInput := referenceframe.FloatsToInputs([]float64{0.5128, -1.9168, -2.2394, -1.8598, 1.1410, -0.9511})

	// End effector frame
	vgPose := spatialmath.NewPoseFromPoint(r3.Vector{Z: 90})
	vgBox, _ := spatialmath.NewBox(spatialmath.NewPoseFromPoint(r3.Vector{Z: 75}), r3.Vector{200, 200, 200}, "")
	vg, _ := referenceframe.NewStaticFrameWithGeometry(testEndEffectorFrame, vgPose, vgBox)

	// Goal pose
	goalPose := spatialmath.NewPose(
		r3.Vector{X: -600.0, Y: -400.0, Z: 60.0},
		&spatialmath.OrientationVectorDegrees{Theta: 0, OX: 0, OY: -1, OZ: 0},
	)

	// Obstacles
	vertWallBehindPose := spatialmath.NewPoseFromPoint(r3.Vector{X: 350.0, Y: 0.0, Z: 500.0})
	mountTablePose := spatialmath.NewPoseFromPoint(r3.Vector{X: 350.0, Y: 0.0, Z: -500.0})
	userTablePose := spatialmath.NewPoseFromPoint(r3.Vector{X: -450.0, Y: 00.0, Z: -266.0})
	worldState, err := referenceframe.WorldStateFromProtobuf(&commonpb.WorldState{
		Obstacles: []*commonpb.GeometriesInFrame{
			{
				ReferenceFrame: "world",
				Geometries: []*commonpb.Geometry{
					{
						Center: spatialmath.PoseToProtobuf(vertWallBehindPose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 80,
								Y: 1000,
								Z: 1000,
							}},
						},
					},
					{
						Center: spatialmath.PoseToProtobuf(mountTablePose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 700,
								Y: 1000,
								Z: 1000,
							}},
						},
					},
					{
						Center: spatialmath.PoseToProtobuf(userTablePose),
						GeometryType: &commonpb.Geometry_Box{
							Box: &commonpb.RectangularPrism{DimsMm: &commonpb.Vector3{
								X: 900,
								Y: 2000,
								Z: 100,
							}},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &config{
		Start:            startInput,
		Goal:             goalPose,
		RobotFrame:       model,
		EndEffectorFrame: &vg,
		WorldState:       worldState,
	}, nil
}

func calcPose(pos []float64) spatialmath.Pose {
	positions := map[string][]referenceframe.Input{}
	inputs := referenceframe.FloatsToInputs(pos)
	positions["arm"] = inputs
	posFrame := referenceframe.NewPoseInFrame("arm", spatialmath.NewZeroPose())
	tf, err := sceneFS.Transform(positions, posFrame, "world")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pose, _ := tf.(*referenceframe.PoseInFrame)
	return pose.Pose()
}
