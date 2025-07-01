/**
 * 农机图片资源
 * 包含了各种农业机械图片的URL地址
 */

// 导入图片资源
import jiqi1 from './jiqi1.png';
import jiqi2 from './jiqi2.png';
import jiqi3 from './jiqi3.png';

// 定义图片对象的接口
export interface CarouselImage {
  url: string | any;  // 修改为接受string或图片资源
  alt: string;
  description?: string;
}

export interface MachineryProduct {
  id: string;
  name: string;
  images: (string | any)[];  // 修改为接受string或图片资源
}

// 收割机轮播图片
export const carouselImages: CarouselImage[] = [
  {
    url: jiqi1, 
    alt: 'CLAAS联合收割机',
    description: '德国CLAAS LEXION系列联合收割机，适用于大规模麦田收割'
  },
  {
    url: jiqi2,
    alt: '约翰迪尔收割机',
    description: '美国约翰迪尔S系列联合收割机，采用先进智能收割技术'
  },
  {
    url: jiqi3,
    alt: '现代收割机',
    description: '现代农机技术的联合收割机，专为中国农田环境设计'
  }
];

// 农机产品图片
export const machineryProducts: MachineryProduct[] = [
  {
    id: 'HN20240517',
    name: '惠农5100联合收割机',
    images: [jiqi1]
  },
  {
    id: 'HN20240518',
    name: '惠农3088手扶式拖拉机',
    images: [jiqi2]
  },
  {
    id: 'HN20240519',
    name: '惠农2560水稻插秧机',
    images: [jiqi3]
  }
]; 