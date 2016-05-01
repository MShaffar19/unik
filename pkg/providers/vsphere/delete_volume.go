package vsphere

import (
	"github.com/emc-advanced-dev/unik/pkg/types"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func (p *VsphereProvider) DeleteVolume(id string, force bool) error {
	volume, err := p.GetVolume(id)
	if err != nil {
		return lxerrors.New("retrieving volume "+id, err)
	}
	if volume.Attachment != "" {
		if force {
			if err := p.DetachVolume(volume.Id); err != nil {
				return lxerrors.New("detaching volume for deletion", err)
			} else {
				return lxerrors.New("volume "+volume.Id+" is attached to instance."+volume.Attachment+", try again with --force or detach volume first", err)
			}
		}
	}
	volumeDir := getVolumeDatastoreDir(volume.Name)
	err = p.getClient().Rmdir(volumeDir)
	if err != nil {
		return lxerrors.New("could not delete volume at path "+ volumeDir, err)
	}
	err = p.state.ModifyVolumes(func(volumes map[string]*types.Volume) error {
		delete(volumes, volume.Id)
		return nil
	})
	if err != nil {
		return lxerrors.New("deleting volume path from state", err)
	}
	err = p.state.Save()
	if err != nil {
		return lxerrors.New("saving image map to state", err)
	}
	return nil
}
