package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/raphadam/gelly"
	"github.com/raphadam/gelly/resolv"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var GameScene = gelly.Scene{
	Name:   "main",
	Layers: []gelly.Layer{&GameUI{}, &Game{}},
}

// var IntroScene = gelly.Scene{
// 	Name:   "intro",
// 	Layers: []gelly.Layer{&Intro{}},
// }

// var PrensentationScene = gelly.Scene{
// 	Name:   "Prensentation",
// 	Layers: []gelly.Layer{&Prensentation{}},
// }

func main() {
	err := gelly.Run(GameScene)
	log.Fatal("client closed", err)
}

type DialogBox struct {
	dur         time.Duration
	text        string
	elapsedTime time.Duration
	runes       []rune
	ui          ebitenui.UI
	face        font.Face
	area        *widget.Text
}

func (d *DialogBox) Init(dur time.Duration, text string) {
	face, _ := loadFont(26)

	*d = DialogBox{
		dur:   dur,
		text:  text,
		runes: []rune(text),
		face:  face,
	}

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(1200, 400),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)),
		)),
	)

	// widget.NewRowLayout(
	// 	// widget.RowLayoutOpts.Spacing(widget.RowLa)
	// 	widget.RowLayoutOpts.
	// )

	d.area = widget.NewText(
		widget.TextOpts.Text("", face, color.White),
		widget.TextOpts.MaxWidth(1240),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),
	)
	root.AddChild(d.area)

	d.ui = ebitenui.UI{
		Container: root,
	}
}

func (d *DialogBox) Update(dt time.Duration) {
	d.elapsedTime += dt
	d.ui.Update()
}

func (d *DialogBox) Draw(r *ebiten.Image) {
	progress := gelly.Clamp(0, 1, float64(d.elapsedTime)/float64(d.dur))
	num := int(float64(len(d.runes)) * progress)
	d.area.Label = string(d.runes[:num])

	d.ui.Draw(r)
}

func (d *DialogBox) Dispose() {}

type Prensentation struct {
	dialog DialogBox
}

func (i *Prensentation) Init(c *gelly.Client) {
	ebiten.SetWindowSize(1280, 720)
	c.SetLayoutSize(1280, 720)

	i.dialog.Init(5*time.Second, "Hello You are in the Recruting game ! In this game you we will try to see if you are matching. Is this text gonna break. I Don't know let's try to see!!")
}

func (i *Prensentation) Message(c *gelly.Client, msg gelly.Message) bool {
	return false
}

func (i *Prensentation) Update(c *gelly.Client, dt time.Duration) {
	i.dialog.Update(dt)
}

func (i *Prensentation) Draw(r *ebiten.Image) {
	i.dialog.Draw(r)
}

func (i *Prensentation) Dispose(c *gelly.Client) {

}

type Intro struct {
	ui ebitenui.UI
}

func (i *Intro) Init(c *gelly.Client) {
	ebiten.SetWindowSize(1280, 720)
	c.SetLayoutSize(1280/2, 720/2)

	// load images for button states: idle, hover, and pressed
	// buttonImage, _ := loadButtonImage()

	// // load button text font
	// face, _ := loadFont(20)

	// root := widget.NewContainer(
	// 	widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
	// 	widget.ContainerOpts.Layout(widget.NewAnchorLayout(
	// 	// widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(30)),
	// 	)),
	// )

	// construct a button
	// button := widget.NewButton(
	// 	// set general widget options
	// 	widget.ButtonOpts.WidgetOpts(
	// 		// instruct the container's anchor layout to center the button both horizontally and vertically
	// 		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
	// 			HorizontalPosition: widget.AnchorLayoutPositionCenter,
	// 			VerticalPosition:   widget.AnchorLayoutPositionCenter,
	// 		}),
	// 	),

	// 	// specify the images to use
	// 	widget.ButtonOpts.Image(buttonImage),

	// 	// specify the button's text, the font face, and the color
	// 	//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
	// 	widget.ButtonOpts.Text("The Recruiting Game", face, &widget.ButtonTextColor{
	// 		Idle:  color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
	// 		Hover: color.NRGBA{0, 255, 128, 255},
	// 	}),
	// 	widget.ButtonOpts.TextProcessBBCode(true),
	// 	// specify that the button's text needs some padding for correct display
	// 	widget.ButtonOpts.TextPadding(widget.Insets{
	// 		Left:   30,
	// 		Right:  30,
	// 		Top:    5,
	// 		Bottom: 5,
	// 	}),

	// 	// add a handler that reacts to clicking the button
	// 	widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
	// 		c.ChangeScene(PrensentationScene)
	// 	}),

	// 	// // add a handler that reacts to entering the button with the cursor
	// 	// widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
	// 	// 	println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
	// 	// }),

	// 	// // add a handler that reacts to moving the cursor on the button
	// 	// widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
	// 	// 	println("cursor moved on button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY, "diffX =", args.DiffX, "diffY =", args.DiffY)
	// 	// }),

	// 	// // add a handler that reacts to exiting the button with the cursor
	// 	// widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
	// 	// 	println("cursor exited button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
	// 	// }),

	// 	// Indicate that this button should not be submitted when enter or space are pressed
	// 	// widget.ButtonOpts.DisableDefaultKeys(),
	// )

	// add the button as a child of the container
	// root.AddChild(button)

	// innerContainer := widget.NewContainer(
	// 	widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
	// 	widget.ContainerOpts.WidgetOpts(
	// 		widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
	// 			HorizontalPosition: widget.AnchorLayoutPositionCenter,
	// 			VerticalPosition:   widget.AnchorLayoutPositionCenter,
	// 			StretchHorizontal:  false,
	// 			StretchVertical:    false,
	// 		}),
	// 		// widget.WidgetOpts.MinSize(100, 100),
	// 	),
	// )
	// root.AddChild(innerContainer)

	// construct the UI
	i.ui = ebitenui.UI{
		// Container: root,
	}
}

func (i *Intro) Message(c *gelly.Client, msg gelly.Message) bool {
	return false
}

func (i *Intro) Update(c *gelly.Client, dt time.Duration) {
	i.ui.Update()
}

func (i *Intro) Draw(r *ebiten.Image) {
	i.ui.Draw(r)

	// op := &text.DrawOptions{}
	// // op.GeoM.Translate(10, 60)

	// text.Draw(r, "The Recruiting Game", &text.GoTextFace{
	// 	Source: Goface,
	// 	Size:   30,
	// }, op)
}

func (i *Intro) Dispose(c *gelly.Client) {

}

// func loadButtonImage() (*widget.ButtonImage, error) {
// 	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

// 	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

// 	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

// 	return &widget.ButtonImage{
// 		Idle:    idle,
// 		Hover:   hover,
// 		Pressed: pressed,
// 	}, nil
// }

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

const (
	WALL_TAG       string = "WALL_TAG"
	SPIKE_TAG      string = "SPIKE_TAG"
	CHECKPOINT_TAG string = "CHECKPOINT_TAG"
	STRAWBERRY_TAG string = "STRAWBERRY_TAG"
)

type EntityType int

const (
	PLAYER EntityType = iota
	SPIKE
	CHECKPOINT
	STRAWBERRY
)

type Physic struct {
	Acceleration gelly.Vector2
	Speed        gelly.Vector2
	Friction     float64
	MaxSpeed     float64
}

type Direction int

const (
	TOP Direction = iota
	RIGHT
	BOTTOM
	LEFT
)

type Components int

const (
	ASPRITE Components = 1 << iota
	SPRITE
	PHYSIC
)

var GRAVITY = gelly.Vector2{X: 0, Y: 600}

type Flags uint64

const (
	IS_DEAD Flags = 1 << iota
	IS_WIN
	IS_ON_GROUND
	IS_SECOND_JUMP
)

type Entity struct {
	Type          EntityType
	Flags         Flags
	Components    Components
	Direction     Direction
	Asprite       gelly.Asprite
	Sprite        gelly.Sprite
	Physic        Physic
	Object        *resolv.Object
	OnGround      *resolv.Object
	SlidingOnWall *resolv.Object
	FacingRight   bool
	IsFinish      bool
}

type Game struct {
	Player1    gelly.Key
	Checkpoint gelly.Key
	Entites    gelly.Pool[Entity]
	Background gelly.Tilemap
	Space      *resolv.Space
	Camera     gelly.Camera
	Face       font.Face
}

func (g *Game) Init(c *gelly.Client) {
	ebiten.SetWindowSize(1280, 720)
	c.SetLayoutSize(1280/2, 720/2)
	ebiten.SetRunnableOnUnfocused(false)

	face, err := loadFont(24)
	if err != nil {
		log.Fatal("unable to load font")
	}
	g.Face = face

	// g.Camera = gelly.NewFollowingCamera(1280, 720, 1280/2, 720/2)
	// g.Camera.Transform.Position.X = -(1280 / 2)
	// g.Camera.Transform.Position.Y = -(720 / 2)
	// g.Camera.Transform.Scale.X = 100
	// g.Camera.Transform.Scale.Y = 100

	g.Background = Background
	g.Space = resolv.NewSpace(700, 400, 16, 16)

	for _, id := range IntGrid {
		switch id.Value {
		case 1:
			g.Space.Add(resolv.NewObject(float64(id.Position[0]), float64(id.Position[1]), 16, 16, WALL_TAG))

		case 2:
			// p := resolv.NewObject(float64(id.Position[0]), float64(id.Position[1]), 16, 16, SPIKE_TAG)
			// g.Entites.Create(Entity{
			// 	Type:       SPIKE,
			// 	Sprite:     SpikeSprite,
			// 	Object:     p,
			// 	Components: SPRITE | PHYSIC,
			// })
			// g.Space.Add(p)

		case 3:
			p := resolv.NewObject(float64(id.Position[0]), float64(id.Position[1]), 32, 32)
			g.Player1 = g.Entites.Create(Entity{
				Type:       PLAYER,
				Asprite:    MaskManAsprite,
				Object:     p,
				Physic:     Physic{MaxSpeed: 4, Friction: 0.5},
				Components: ASPRITE | PHYSIC,
			})
			g.Space.Add(p)

		case 4:
			p := resolv.NewObject(float64(id.Position[0]), float64(id.Position[1]), 64, 64, CHECKPOINT_TAG)
			g.Checkpoint = g.Entites.Create(Entity{
				Type:       CHECKPOINT,
				Asprite:    CheckpointAsprite,
				Object:     p,
				Components: ASPRITE | PHYSIC,
			})
			g.Space.Add(p)

		case 5:
			p := resolv.NewObject(float64(id.Position[0]), float64(id.Position[1]), 16, 16, STRAWBERRY_TAG)

			key := g.Entites.Create(Entity{
				Type:       STRAWBERRY,
				Asprite:    StrawberriesAsprite,
				Object:     p,
				Components: ASPRITE | PHYSIC,
			})
			p.Key = key

			g.Space.Add(p)
		}

	}

	player, ok := g.Entites.Get(g.Player1)
	if ok {
		player.FacingRight = true
	}
}

func (g *Game) Message(c *gelly.Client, msg gelly.Message) bool {
	return false
}

func (g *Game) Update(c *gelly.Client, dt time.Duration) {
	dt = 16 * time.Millisecond

	// g.Dialog.Update(dt)
	g.Entites.Apply()

	g.Entites.For(func(i int, k gelly.Key, v *Entity) bool {
		if v.Type == PLAYER {
			// update Physic
			{
				ResolvePhysicMovement(v)

				if collision := v.Object.Check(0, 0, CHECKPOINT_TAG); collision != nil {

					checkpoint, ok := g.Entites.Get(g.Checkpoint)
					if ok {

						if !checkpoint.IsFinish {
							WinPlayer.Rewind()
							WinPlayer.Play()
							checkpoint.IsFinish = true
							checkpoint.Asprite.Change("finish")
						}
					}
				}

				if collision := v.Object.Check(0, 0, STRAWBERRY_TAG); collision != nil {
					stawberry, ok := g.Entites.Get(collision.Objects[0].Key)
					if ok {
						PickUpPlayer.Rewind()
						PickUpPlayer.Play()
						g.Space.Remove(stawberry.Object)
						g.Entites.Destroy(collision.Objects[0].Key)
					}
				}

				// if collision := v.Object.Check(movement.X, movement.Y, SPIKE_TAG); collision != nil {
				// 	v.Asprite.Change("doubleJump")
				// 	v.Flags |= IS_DEAD
				// }

				// if collision := v.Object.Check(movement.X, movement.Y, CHECKPOINT_TAG); collision != nil {
				// 	log.Println("WINNING!!!")
				// 	// g.Entites.Destroy(k)
				// 	v.Flags |= IS_WIN
				// }

				// v.Object.Position = v.Object.Position.Add(movement)
				// v.Object.Update()
			}

			// // update death if necessary
			// {
			// 	if v.Flags&IS_DEAD != 0 || v.Flags&IS_WIN != 0 {
			// 		g.Space.Remove(v.Object)
			// 		g.Entites.Destroy(k)
			// 		log.Println("should be dying")
			// 	}
			// }

			// update asprite
			{
				abs := math.Abs(v.Physic.Speed.X)
				if abs > 0 {
					v.Asprite.Change("walk")
				} else {
					v.Asprite.Change("idle")
				}

				v.Asprite.FlipH = !v.FacingRight

				if v.Physic.Speed.Y > 0.1 {
					v.Asprite.Change("fall")
				}
			}
		}

		if v.Components&ASPRITE > 0 {
			v.Asprite.Update(dt)
		}

		return false
	})
}

func ResolvePhysicMovement(player *Entity) {
	// Floating platform movement needs to be done before the player's movement update to make sure there's no space between its top and the player's bottom;
	// otherwise, an alternative might be to have the platform detect to see if the Player's resting on it, and if so, move the player up manually.
	// y, _, seqDone := world.FloatingPlatformTween.Update(1.0 / 60.0)
	// world.FloatingPlatform.Position.Y = float64(y)
	// if seqDone {
	// 	world.FloatingPlatformTween.Reset()
	// }
	// world.FloatingPlatform.Update()

	// Now we update the Player's movement. This is the real bread-and-butter of this example, naturally.
	friction := 0.4
	accel := 0.4 + friction
	maxSpeed := 4.0
	jumpSpd := 10.0
	gravity := 0.65

	player.Physic.Speed.Y += gravity

	// if player.SlidingOnWall != nil && player.Physic.Speed.Y > 1 {
	// 	player.Physic.Speed.Y = 1
	// }

	// Horizontal movement is only possible when not wallsliding.
	if player.SlidingOnWall == nil {
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxisValue(0, 0) > 0.1 {
			player.Physic.Speed.X += accel
			player.FacingRight = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxisValue(0, 0) < -0.1 {
			player.Physic.Speed.X -= accel
			player.FacingRight = false
		}
	}

	// Apply friction and horizontal speed limiting.
	if player.Physic.Speed.X > friction {
		player.Physic.Speed.X -= friction
	} else if player.Physic.Speed.X < -friction {
		player.Physic.Speed.X += friction
	} else {
		player.Physic.Speed.X = 0
	}

	if player.Physic.Speed.X > maxSpeed {
		player.Physic.Speed.X = maxSpeed
	} else if player.Physic.Speed.X < -maxSpeed {
		player.Physic.Speed.X = -maxSpeed
	}

	// Check for jumping.
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsGamepadButtonJustPressed(0, 0) {

		// if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.GamepadAxisValue(0, 1) > 0.1) && player.OnGround != nil && player.OnGround.HasTags("platform") {

		// 	// player.IgnorePlatform = player.OnGround

		// } else {

		if player.OnGround != nil {
			player.Physic.Speed.Y = -jumpSpd
			JumpPlayer.Rewind()
			JumpPlayer.Play()

		} else if player.SlidingOnWall != nil {
			// WALLJUMPING
			player.Physic.Speed.Y = -jumpSpd

			if player.SlidingOnWall.Position.X > player.Object.Position.X {
				player.Physic.Speed.X = -4
			} else {
				player.Physic.Speed.X = 4
			}

			player.SlidingOnWall = nil
		}

		// }
	}

	// We handle horizontal movement separately from vertical movement. This is, conceptually, decomposing movement into two phases / axes.
	// By decomposing movement in this manner, we can handle each case properly (i.e. stop movement horizontally separately from vertical movement, as
	// necesseary). More can be seen on this topic over on this blog post on higherorderfun.com:
	// http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

	// dx is the horizontal delta movement variable (which is the Player's horizontal speed). If we come into contact with something, then it will
	// be that movement instead.
	dx := player.Physic.Speed.X

	// Moving horizontally is done fairly simply; we just check to see if something solid is in front of us. If so, we move into contact with it
	// and stop horizontal movement speed. If not, then we can just move forward.

	if check := player.Object.Check(player.Physic.Speed.X, 0, "WALL_TAG"); check != nil {

		dx = check.ContactWithCell(check.Cells[0]).X
		player.Physic.Speed.X = 0

		// If you're in the air, then colliding with a wall object makes you start wall sliding.
		// if player.OnGround == nil {
		// 	player.SlidingOnWall = check.Objects[0]
		// }

	}

	// Then we just apply the horizontal movement to the Player's Object. Easy-peasy.
	player.Object.Position.X += dx

	// Now for the vertical movement; it's the most complicated because we can land on different types of objects and need
	// to treat them all differently, but overall, it's not bad.

	// First, we set OnGround to be nil, in case we don't end up standing on anything.
	player.OnGround = nil

	// dy is the delta movement downward, and is the vertical movement by default; similarly to dx, if we come into contact with
	// something, this will be changed to move to contact instead.

	dy := player.Physic.Speed.Y

	// We want to be sure to lock vertical movement to a maximum of the size of the Cells within the Space
	// so we don't miss any collisions by tunneling through.

	dy = math.Max(math.Min(dy, 16), -16)

	// We're going to check for collision using dy (which is vertical movement speed), but add one when moving downwards to look a bit deeper down
	// into the ground for solid objects to land on, specifically.
	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	// We check for any solid / stand-able objects. In actuality, there aren't any other Objects
	// with other tags in this Space, so we don't -have- to specify any tags, but it's good to be specific for clarity in this example.
	if check := player.Object.Check(0, checkDistance, "WALL_TAG", "platform", "ramp"); check != nil {

		// So! Firstly, we want to see if we jumped up into something that we can slide around horizontally to avoid bumping the Player's head.

		// Sliding around a misspaced jump is a small thing that makes jumping a bit more forgiving, and is something different polished platformers
		// (like the 2D Mario games) do to make it a smidge more comfortable to play. For a visual example of this, see this excellent devlog post
		// from the extremely impressive indie game, Leilani's Island: https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138

		// To accomplish this sliding, we simply call Collision.SlideAgainstCell() to see if we can slide.
		// We pass the first cell, and tags that we want to avoid when sliding (i.e. we don't want to slide into cells that contain other solid objects).

		// slide, slideOK := check.SlideAgainstCell(check.Cells[0], "solid")

		// We further ensure that we only slide if:
		// 1) We're jumping up into something (dy < 0),
		// 2) If the cell we're bumping up against contains a solid object,
		// 3) If there was, indeed, a valid slide left or right, and
		// 4) If the proposed slide is less than 8 pixels in horizontal distance. (This is a relatively arbitrary number that just so happens to be half the
		// width of a cell. This is to ensure the player doesn't slide too far horizontally.)

		// if dy < 0 && check.Cells[0].ContainsTags("solid") && slideOK && math.Abs(slide.X) <= 8 {

		// If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
		// player.Object.Position.X += slide.X
		// } else

		{

			// If sliding -fails-, that means the Player is jumping directly onto or into something, and we need to do more to see if we need to come into
			// contact with it. Let's press on!

			// First, we check for ramps. For ramps, we can't simply check for collision with Check(), as that's not precise enough. We need to get a bit
			// more information, and so will do so by checking its Shape (a triangular ConvexPolygon, as defined in WorldPlatformer.Init()) against the
			// Player's Shape (which is also a rectangular ConvexPolygon).

			// We get the ramp by simply filtering out Objects with the "ramp" tag out of the objects returned in our broad Check(), and grabbing the first one
			// if there's any at all.
			// if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {

			// 	// For simplicity, this code assumes we can only stand on one ramp at a time as there is only one ramp in this example.
			// 	// This is exemplified by the ramp := ramps[0] line.
			// 	// In actuality, if there was a possibility to have a potential collision with multiple ramps (i.e. a ramp that sits on another ramp, and the player running down
			// 	// one onto the other), the collision testing code should probably go with the ramp with the highest confirmed intersection point out of the two.

			// 	ramp := ramps[0]

			// 	// Next, we see if there's been an intersection between the two Shapes using Shape.Intersection. We pass the ramp's shape, and also the movement
			// 	// we're trying to make horizontally, as this makes the Intersection function return the next y-position while moving, not the one directly
			// 	// underneath the Player. This would keep the player from getting "stuck" when walking up a ramp into the top of a solid block, if there weren't
			// 	// a landing at the top and bottom of the ramp.

			// 	// We use 8 here for the Y-delta so that we can easily see if you're running down the ramp (in which case you're probably in the air as you
			// 	// move faster than you can fall in this example). This way we can maintain contact so you can always jump while running down a ramp. We only
			// 	// continue with coming into contact with the ramp as long as you're not moving upwards (i.e. jumping).

			// 	if contactSet := player.Object.Shape.Intersection(dx, 8, ramp.Shape); dy >= 0 && contactSet != nil {

			// 		// If Intersection() is successful, a ContactSet is returned. A ContactSet contains information regarding where
			// 		// two Shapes intersect, like the individual points of contact, the center of the contacts, and the MTV, or
			// 		// Minimum Translation Vector, to move out of contact.

			// 		// Here, we use ContactSet.TopmostPoint() to get the top-most contact point as an indicator of where
			// 		// we want the player's feet to be. Then we just set that position with a tiny bit of collision margin,
			// 		// and we're done.

			// 		dy = contactSet.TopmostPoint().Y - player.Object.Bottom() + 0.1
			// 		player.OnGround = ramp
			// 		player.Physic.Speed.Y = 0

			// 	}

			// }

			// Platforms are next; here, we just see if the platform is not being ignored by attempting to drop down,
			// if the player is falling on the platform (as otherwise he would be jumping through platforms), and if the platform is low enough
			// to land on. If so, we stand on it.

			// Because there's a moving floating platform, we use Collision.ContactWithObject() to ensure the player comes into contact
			// with the top of the platform object. An alternative would be to use Collision.ContactWithCell(), but that would be only if the
			// platform didn't move and were aligned with the Space's grid.

			// if platforms := check.ObjectsByTags("platform"); len(platforms) > 0 {

			// 	platform := platforms[0]

			// 	if platform != player.IgnorePlatform && dy >= 0 && player.Object.Bottom() < platform.Position.Y+4 {
			// 		dy = check.ContactWithObject(platform).Y
			// 		player.OnGround = platform
			// 		player.Physic.Speed.Y = 0
			// 	}

			// }

			// Finally, we check for simple solid ground. If we haven't had any success in landing previously, or the solid ground
			// is higher than the existing ground (like if the platform passes underneath the ground, or we're walking off of solid ground
			// onto a ramp), we stand on it instead. We don't check for solid collision first because we want any ramps to override solid
			// ground (so that you can walk onto the ramp, rather than sticking to solid ground).

			// We use ContactWithObject() here because otherwise, we might come into contact with the moving platform's cells (which, naturally,
			// would be selected by a Collision.ContactWithCell() call because the cell is closest to the Player).

			if solids := check.ObjectsByTags("WALL_TAG"); len(solids) > 0 && (player.OnGround == nil || player.OnGround.Position.Y >= solids[0].Position.Y) {
				dy = check.ContactWithObject(solids[0]).Y
				player.Physic.Speed.Y = 0

				// We're only on the ground if we land on it (if the object's Y is greater than the player's).
				if solids[0].Position.Y > player.Object.Position.Y {
					player.OnGround = solids[0]
				}

			}

			// if player.OnGround != nil {
			// 	// player.SlidingOnWall = nil  // Player's on the ground, so no wallsliding anymore.
			// 	// player.IgnorePlatform = nil // Player's on the ground, so reset which platform is being ignored.
			// }

		}

	}

	// Move the object on dy.
	player.Object.Position.Y += dy

	// wallNext := 1.0
	// if !player.FacingRight {
	// 	wallNext = -1
	// }

	// If the wall next to the Player runs out, stop wall sliding.
	// if c := player.Object.Check(wallNext, 0, "solid"); player.SlidingOnWall != nil && c == nil {
	// 	player.SlidingOnWall = nil
	// }

	player.Object.Update() // Update the player's position in the space.

	// And that's it!
}

func (g *Game) Draw(r *ebiten.Image) {
	// if player, ok := g.Entites.Get(g.Player1); ok {
	// 	v := gelly.Vector2(player.Object.Center())
	// 	v.Y -= 30
	// 	g.Camera.Follow(v)
	// }

	g.Background.Draw(r)

	text.Draw(r, "You", g.Face, 105, 270, gelly.ColorWhite)
	text.Draw(r, "Goal", g.Face, 517, 270, gelly.ColorWhite)

	gelly.Sprite{
		Image: KeyboardArrowLeftImg,
		Transform: gelly.Transform{
			Position: gelly.Vector2{X: 200, Y: 50},
			Scale:    gelly.Vector2{X: -40, Y: -40},
		},
	}.Draw(r)
	text.Draw(r, "Left", g.Face, 200, 120, gelly.ColorWhite)
	text.Draw(r, "Gauche", g.Face, 200, 150, gelly.ColorWhite)

	gelly.Sprite{
		Image: KeyboardArrowUpImg,
		Transform: gelly.Transform{
			Position: gelly.Vector2{X: 300, Y: 50},
			Scale:    gelly.Vector2{X: -40, Y: -40},
		},
	}.Draw(r)
	text.Draw(r, "Jump", g.Face, 300, 120, gelly.ColorWhite)
	text.Draw(r, "Sauter", g.Face, 300, 150, gelly.ColorWhite)

	gelly.Sprite{
		Image: KeyboardArrowRightImg,
		Transform: gelly.Transform{
			Position: gelly.Vector2{X: 400, Y: 50},
			Scale:    gelly.Vector2{X: -40, Y: -40},
		},
	}.Draw(r)
	text.Draw(r, "Right", g.Face, 400, 120, gelly.ColorWhite)
	text.Draw(r, "Droite", g.Face, 400, 150, gelly.ColorWhite)

	g.Entites.For(func(i int, k gelly.Key, e *Entity) bool {

		if e.Components&ASPRITE > 0 {
			e.Asprite.Transform.Position = gelly.Vector2(e.Object.Position)
			e.Asprite.Draw(r)
		}

		if e.Components&SPRITE > 0 {
			e.Sprite.Transform.Position = gelly.Vector2(e.Object.Position)
			e.Sprite.Draw(r)
		}

		return false
	})

	// debug physic
	// for _, obj := range g.Space.Objects() {
	// 	vector.StrokeRect(r,
	// 		float32(obj.Position.X),
	// 		float32(obj.Position.Y),
	// 		float32(obj.Size.X),
	// 		float32(obj.Size.Y),
	// 		1, gelly.ColorCyan, false,
	// 	)
	// }

	// g.Camera.Draw(r)
	// g.Dialog.Draw(r)
}

func (g *Game) Dispose(c *gelly.Client) {
}

type GameUI struct {
}

func (g *GameUI) Init(c *gelly.Client) {
	log.Println("On the new level")
}

func (g *GameUI) Message(c *gelly.Client, msg gelly.Message) bool {
	return false
}

func (g *GameUI) Update(c *gelly.Client, dt time.Duration) {

}

func (g *GameUI) Draw(r *ebiten.Image) {

}

func (g *GameUI) Dispose(c *gelly.Client) {

}
