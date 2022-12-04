package navigator

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ikawaha/navigator/event"
	"github.com/ikawaha/navigator/webdriver/session"
)

type actionsFunc func(*session.Element) error

func (s *Selection) forEachElement(ctx context.Context, actions actionsFunc) error {
	elements, err := s.getElementsAtLeastOne(ctx)
	if err != nil {
		return fmt.Errorf("failed to select elements from %s: %w", s, err)
	}
	for _, element := range elements {
		if err := actions(element); err != nil {
			return err
		}
	}
	return nil
}

// Click clicks on all the elements that the selection refers to.
func (s *Selection) Click() error {
	return s.ClickWithContext(context.Background())
}

// ClickWithContext clicks on all the elements that the selection refers to.
func (s *Selection) ClickWithContext(ctx context.Context) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := selectedElement.Click(ctx); err != nil {
			return fmt.Errorf("failed to click on %s: %w", s, err)
		}
		return nil
	})
}

// DoubleClick double-clicks on all the elements that the selection refers to.
func (s *Selection) DoubleClick() error {
	return s.DoubleClickWithContext(context.Background())
}

// DoubleClickWithContext double-clicks on all the elements that the selection refers to.
func (s *Selection) DoubleClickWithContext(ctx context.Context) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := s.session.MoveTo(ctx, selectedElement, nil); err != nil {
			return fmt.Errorf("failed to move mouse to %s: %w", s, err)
		}
		if err := s.session.DoubleClick(ctx); err != nil {
			return fmt.Errorf("failed to double-click on %s: %w", s, err)
		}
		return nil
	})
}

// Clear clears all fields the selection refers to.
func (s *Selection) Clear() error {
	return s.ClearWithContext(context.Background())
}

// ClearWithContext clears all fields the selection refers to.
func (s *Selection) ClearWithContext(ctx context.Context) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := selectedElement.Clear(ctx); err != nil {
			return fmt.Errorf("failed to clear %s: %w", s, err)
		}
		return nil
	})
}

// Fill fills all the fields the selection refers to with the provided text.
func (s *Selection) Fill(text string) error {
	return s.FillWithContext(context.Background(), text)
}

// FillWithContext fills all the fields the selection refers to with the provided text.
func (s *Selection) FillWithContext(ctx context.Context, text string) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := selectedElement.Clear(ctx); err != nil {
			return fmt.Errorf("failed to clear %s: %w", s, err)
		}
		if err := selectedElement.Value(ctx, text); err != nil {
			return fmt.Errorf("failed to enter text into %s: %w", s, err)
		}
		return nil
	})
}

// UploadFile uploads the provided file to all selected <input type="file" />.
// The provided filename may be a relative or absolute path.
// Returns an error if elements of any other type are in the selection.
func (s *Selection) UploadFile(filename string) error {
	return s.UploadFileWithContext(context.Background(), filename)
}

// UploadFileWithContext uploads the provided file to all selected <input type="file" />.
// The provided filename may be a relative or absolute path.
// Returns an error if elements of any other type are in the selection.
func (s *Selection) UploadFileWithContext(ctx context.Context, filename string) error {
	absFilePath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("failed to find absolute path for filename: %w", err)
	}
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		tagName, err := selectedElement.GetName(ctx)
		if err != nil {
			return fmt.Errorf("failed to determine tag name of %s: %w", s, err)
		}
		if tagName != "input" {
			return fmt.Errorf("element for %s is not an input element", s)
		}
		inputType, err := selectedElement.GetAttribute(ctx, "type")
		if err != nil {
			return fmt.Errorf("failed to determine type attribute of %s: %w", s, err)
		}
		if inputType != "file" {
			return fmt.Errorf("element for %s is not a file uploader", s)
		}
		if err := selectedElement.Value(ctx, absFilePath); err != nil {
			return fmt.Errorf("failed to enter text into %s: %w", s, err)
		}
		return nil
	})
}

// Check checks all the unchecked checkboxes that the selection refers to.
func (s *Selection) Check() error {
	return s.CheckWithContext(context.Background())
}

// CheckWithContext checks all the unchecked checkboxes that the selection refers to.
func (s *Selection) CheckWithContext(ctx context.Context) error {
	return s.setChecked(ctx, true)
}

// Uncheck unchecks all the checked checkboxes that the selection refers to.
func (s *Selection) Uncheck() error {
	return s.UncheckWithContext(context.Background())
}

// UncheckWithContext unchecks all the checked checkboxes that the selection refers to.
func (s *Selection) UncheckWithContext(ctx context.Context) error {
	return s.setChecked(ctx, false)
}

func (s *Selection) setChecked(ctx context.Context, checked bool) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		elementType, err := selectedElement.GetAttribute(ctx, "type")
		if err != nil {
			return fmt.Errorf("failed to retrieve type attribute of %s: %w", s, err)
		}
		if elementType != "checkbox" {
			return fmt.Errorf("%s does not refer to a checkbox", s)
		}
		elementChecked, err := selectedElement.IsSelected(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve state of %s: %w", s, err)
		}
		if elementChecked != checked {
			if err := selectedElement.Click(ctx); err != nil {
				return fmt.Errorf("failed to click on %s: %w", s, err)
			}
		}
		return nil
	})
}

// Select may be called on a selection of any number of <select> elements to select
// any <option> elements under those <select> elements that match the provided text.
func (s *Selection) Select(text string) error {
	return s.SelectWithContext(context.Background(), text)
}

// SelectWithContext may be called on a selection of any number of <select> elements to select
// any <option> elements under those <select> elements that match the provided text.
func (s *Selection) SelectWithContext(ctx context.Context, text string) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		optionXPath := fmt.Sprintf(`./option[normalize-space()="%s"]`, text)
		optionToSelect := selector{Type: xPathType, Value: optionXPath}
		options, err := selectedElement.Session.GetElements(ctx, optionToSelect.SessionSelector())
		if err != nil {
			return fmt.Errorf("failed to select specified option for %s: %w", s, err)
		}
		if len(options) == 0 {
			return fmt.Errorf(`no options with text "%s" found for %s`, text, s)
		}
		for _, option := range options {
			if err := option.Click(ctx); err != nil {
				return fmt.Errorf(`failed to click on option with text "%s" for %s: %w`, text, s, err)
			}
		}
		return nil
	})
}

// Submit submits all selected forms. The selection may refer to a form itself
// or any input element contained within a form.
func (s *Selection) Submit() error {
	return s.SubmitWithContext(context.Background())
}

// SubmitWithContext submits all selected forms. The selection may refer to a form itself
// or any input element contained within a form.
func (s *Selection) SubmitWithContext(ctx context.Context) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := selectedElement.Submit(ctx); err != nil {
			return fmt.Errorf("failed to submit %s: %w", s, err)
		}
		return nil
	})
}

// Tap performs the provided Tap event on each element in the selection.
func (s *Selection) Tap(tap event.Tap) error {
	return s.TapWithContext(context.Background(), tap)
}

// TapWithContext performs the provided Tap event on each element in the selection.
func (s *Selection) TapWithContext(ctx context.Context, tap event.Tap) error {
	var touchFunc func(context.Context, *session.Element) error
	switch tap {
	case event.SingleTap:
		touchFunc = s.session.TouchClick
	case event.DoubleTap:
		touchFunc = s.session.TouchDoubleClick
	case event.LongTap:
		touchFunc = s.session.TouchLongClick
	default:
		return fmt.Errorf("failed to %s on %s: invalid tap event", tap, s)
	}

	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := touchFunc(ctx, selectedElement); err != nil {
			return fmt.Errorf("failed to %s on %s: %w", tap, s, err)
		}
		return nil
	})
}

// Touch performs the provided Touch event at the location of each element in the
// selection.
func (s *Selection) Touch(touch event.Touch) error {
	return s.TouchWithContext(context.Background(), touch)
}

// TouchWithContext performs the provided Touch event at the location of each element in the
// selection.
func (s *Selection) TouchWithContext(ctx context.Context, touch event.Touch) error {
	var touchFunc func(context.Context, int, int) error
	switch touch {
	case event.HoldFinger:
		touchFunc = s.session.TouchDown
	case event.ReleaseFinger:
		touchFunc = s.session.TouchUp
	case event.MoveFinger:
		touchFunc = s.session.TouchMove
	default:
		return fmt.Errorf("failed to %s on %s: invalid touch event", touch, s)
	}

	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		x, y, err := selectedElement.GetLocation(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve location of %s: %w", s, err)
		}
		if err := touchFunc(ctx, x, y); err != nil {
			return fmt.Errorf("failed to flick finger on %s: %w", s, err)
		}
		return nil
	})
}

// FlickFinger performs a flick touch action by the provided offset and at the
// provided speed on exactly one element.
func (s *Selection) FlickFinger(xOffset, yOffset int, speed uint) error {
	return s.FlickFingerWithContext(context.Background(), xOffset, yOffset, speed)
}

// FlickFingerWithContext performs a flick touch action by the provided offset and at the
// provided speed on exactly one element.
func (s *Selection) FlickFingerWithContext(ctx context.Context, xOffset, yOffset int, speed uint) error {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	if err := s.session.TouchFlick(ctx, selectedElement, session.XYOffset{X: xOffset, Y: yOffset}, session.ScalarSpeed(speed)); err != nil {
		return fmt.Errorf("failed to flick finger on %s: %w", s, err)
	}
	return nil
}

// ScrollFinger performs a scroll touch action by the provided offset on exactly
// one element.
func (s *Selection) ScrollFinger(xOffset, yOffset int) error {
	return s.ScrollFingerWithContext(context.Background(), xOffset, yOffset)
}

// ScrollFingerWithContext performs a scroll touch action by the provided offset on exactly
// one element.
func (s *Selection) ScrollFingerWithContext(ctx context.Context, xOffset, yOffset int) error {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	if err := s.session.TouchScroll(ctx, selectedElement, session.XYOffset{X: xOffset, Y: yOffset}); err != nil {
		return fmt.Errorf("failed to scroll finger on %s: %w", s, err)
	}
	return nil
}

// SendKeys sends key events to the selected elements.
func (s *Selection) SendKeys(key string) error {
	return s.SendKeysWithContext(context.Background(), key)
}

// SendKeysWithContext sends key events to the selected elements.
func (s *Selection) SendKeysWithContext(ctx context.Context, key string) error {
	return s.forEachElement(ctx, func(selectedElement *session.Element) error {
		if err := selectedElement.Value(ctx, key); err != nil {
			return fmt.Errorf("failed to send key %s on %s: %w", key, s, err)
		}
		return nil
	})
}

// SwitchToFrame focuses on the frame specified by the selection. All new and
// existing selections will refer to the new frame. All further Page methods
// will apply to this frame as well.
func (s *Selection) SwitchToFrame() error {
	return s.SwitchToFrameWithContext(context.Background())
}

// SwitchToFrameWithContext focuses on the frame specified by the selection. All new and
// existing selections will refer to the new frame. All further Page methods
// will apply to this frame as well.
func (s *Selection) SwitchToFrameWithContext(ctx context.Context) error {
	selectedElement, err := s.getElementExactlyOne(ctx)
	if err != nil {
		return fmt.Errorf("failed to select element from %s: %w", s, err)
	}
	if err := s.session.Frame(ctx, selectedElement); err != nil {
		return fmt.Errorf("failed to switch to frame referred to by %s: %w", s, err)
	}
	return nil
}
