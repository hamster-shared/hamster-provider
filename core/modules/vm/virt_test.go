package vm

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCreate(t *testing.T) {
	cli, err := NewVirtManager()

	if err != nil {
		t.Log(err)
		t.Error(err)
	}

	template := Template{
		Cpu:    1,
		Memory: 1,
		Dist:   50,
		System: "ubuntu",
		Name:   "test",
	}
	cli.SetTemplate(template)

	err = cli.Create()
	if err != nil {
		t.Log(err)
		t.Error(err)
	}
}

func TestMemory(t *testing.T) {

	fmt.Printf("num : %s", fmt.Sprint(1<<10))
}

func TestPublicKey(t *testing.T) {
	cli, err := NewVirtManager()
	assert.NoError(t, err)
	template := Template{
		Cpu:    1,
		Memory: 1,
		Dist:   50,
		System: "ubuntu",
		Name:   "test",
	}
	cli.SetTemplate(template)
	err = cli.InjectionPublicKey("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCWCDS4Io8+PFqGepqy0YNrtw3B7g7lhg7WNcH2VyJmmHlvft69N3S4EzDugEUDgPbgihiL56wyq56GtOG6+RuRuqkEU983MRC6j0yazem/KPs2nAS0NW5A8Nzxm9ixXnF9Bw6qHpO+L8ZbKdsIR+xux5QVriWTmDd/FeaovzRa/Ogr/BdShsp5H1s8aKkj2ygm16rlWAuQcoPQJPDWJVLM9cub8wj/AGrOzRDQCMnbcm69BZT7GPbodVmBIlugICuSVVvKSpZEa0QHCdQW2z2kIan7EwEI7LPYyDpCRAAI2mYEsl9WIIzae1ACK7dKwp9DKfLlKU4YRNfvR5stGgNezelz2pbN0TvK0T6NrqlKDo1eZQbHzRzvUKtDCiwSBdauJVus5Zowqy8lXr9wVosbra8z8cd+vM+e5+82fEjnE3BQm6NUHatOfxe/1MtYeem1Zlru5ISc25ceCXJd/l6qUrlIamHKgMpxvAt8g9pcPpPH2YozLkRlohcdWrhA+kk= gr@gr-Lenovo")
	assert.NoError(t, err)
}

func TestGetIpAddress(t *testing.T) {
	cli, err := NewVirtManager()
	assert.NoError(t, err)
	template := Template{
		Cpu:    1,
		Memory: 1,
		Dist:   50,
		System: "ubuntu",
		Name:   "order_5",
	}
	cli.SetTemplate(template)
	ip, err := cli.GetIp()
	assert.NoError(t, err)
	fmt.Println(ip)
}

func TestGetStatus(t *testing.T) {
	cli, err := NewVirtManager()
	assert.NoError(t, err)
	template := Template{
		Cpu:    1,
		Memory: 1,
		Dist:   50,
		System: "ubuntu",
		Name:   "order_4",
	}
	cli.SetTemplate(template)
	status, err := cli.Status()
	assert.NotNil(t, status)
	assert.NoError(t, err)
}

func TestName(t *testing.T) {
	str := "/home/gr/app/CentOS7.qcow2.bak2"
	order := "order1"
	newPath := fmt.Sprintf("%s/%s_%s", path.Dir(str), order, path.Base(str))

	fmt.Println(newPath)

}

func TestDownload(t *testing.T) {
	//downloadUrl := "https://s3.ttchain.tntlinking.com/compute/CentOS7.qcow2.tar.gz"
	dest := "/tmp/CentOS7.qcow2.tar.gz"
	//err := utils.Download(downloadUrl, dest)
	//assert.NoError(t, err)

	file, err := os.Open(dest)
	err = utils.UnTar(file, "/tmp")
	assert.NoError(t, err)
}
