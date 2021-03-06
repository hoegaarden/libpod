package image

import "github.com/pkg/errors"

// GetPruneImages returns a slice of images that have no names/unused
func (ir *Runtime) GetPruneImages(all bool) ([]*Image, error) {
	var (
		pruneImages []*Image
	)
	allImages, err := ir.GetImages()
	if err != nil {
		return nil, err
	}
	for _, i := range allImages {
		if len(i.Names()) == 0 {
			pruneImages = append(pruneImages, i)
			continue
		}
		if all {
			containers, err := i.Containers()
			if err != nil {
				return nil, err
			}
			if len(containers) < 1 {
				pruneImages = append(pruneImages, i)
			}
		}
	}
	return pruneImages, nil
}

// PruneImages prunes dangling and optionally all unused images from the local
// image store
func (ir *Runtime) PruneImages(all bool) ([]string, error) {
	var prunedCids []string
	pruneImages, err := ir.GetPruneImages(all)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get images to prune")
	}
	for _, p := range pruneImages {
		if err := p.Remove(true); err != nil {
			return nil, errors.Wrap(err, "failed to prune image")
		}
		prunedCids = append(prunedCids, p.ID())
	}
	return prunedCids, nil
}
