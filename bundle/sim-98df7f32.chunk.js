import{dY as t,L as e,D as n,b as a,fs as i,ft as s,aq as r,E as d,co as o,cj as c,cm as l,aD as m,bn as h,B as f,w as p,aE as g,bD as u,A as S,ar as O,a2 as y,F as T}from"./detailed_results-7b150079.chunk.js";import{B as P,m as b,I,T as v}from"./individual_sim_ui-7ca50b32.chunk.js";const w={name:"Celestial Focus",data:t.create({talentsString:"05320031103--230023312131502331050313051",glyphs:e.create({major1:n.GlyphOfWildGrowth,major2:n.GlyphOfSwiftmend,major3:n.GlyphOfNourish,minor2:a.GlyphOfUnburdenedRebirth,minor3:a.GlyphOfTheWild,minor1:a.GlyphOfDash})})},G={name:"Thicc Resto",data:t.create({talentsString:"05320001--230023312331502531053313051",glyphs:e.create({major1:n.GlyphOfWildGrowth,major2:n.GlyphOfSwiftmend,major3:n.GlyphOfNourish,minor2:a.GlyphOfUnburdenedRebirth,minor3:a.GlyphOfTheWild,minor1:a.GlyphOfDash})})},E=i.create({}),F=s.create({innervateTarget:r.create()}),W=d.create({defaultPotion:o.RunicManaPotion,flask:c.FlaskOfTheFrostWyrm,food:l.FoodFishFeast}),j=m.create({arcaneBrilliance:!0,bloodlust:!0,divineSpirit:!0,giftOfTheWild:h.TristateEffectImproved,icyTalons:!0,moonkinAura:h.TristateEffectImproved,leaderOfThePack:h.TristateEffectImproved,powerWordFortitude:h.TristateEffectImproved,sanctifiedRetribution:!0,strengthOfEarthTotem:h.TristateEffectImproved,trueshotAura:!0,wrathOfAirTotem:!0}),k=f.create({blessingOfKings:!0,blessingOfMight:h.TristateEffectImproved,blessingOfWisdom:h.TristateEffectImproved,vampiricTouch:!0}),C=p.create({heroicPresence:!1}),B=g.create({bloodFrenzy:!0,ebonPlaguebringer:!0,faerieFire:h.TristateEffectImproved,heartOfTheCrusader:!0,judgementOfWisdom:!0,shadowMastery:!0,sunderArmor:!0,totemOfWrath:!0}),D={distanceFromTarget:18},M={name:"Pre-raid Preset",tooltip:P,gear:u.fromJsonString('{ "items": [\n\t\t{"id":37149,"enchant":3819,"gems":[41401,40051]},\n\t\t{"id":42339,"gems":[40026]},\n\t\t{"id":37673,"enchant":3809,"gems":[39998]},\n\t\t{"id":41610,"enchant":3831},\n\t\t{"id":42102,"enchant":3832},\n\t\t{"id":37361,"enchant":2332,"gems":[0]},\n\t\t{"id":42113,"enchant":3246,"gems":[0]},\n\t\t{"id":37643,"enchant":3601,"gems":[39998]},\n\t\t{"id":37791,"enchant":3719},\n\t\t{"id":44202,"enchant":3232,"gems":[39998]},\n\t\t{"id":37694},\n\t\t{"id":37192},\n\t\t{"id":37111},\n\t\t{"id":37657},\n\t\t{"id":37169,"enchant":3834},\n\t\t{"id":40699},\n\t\t{"id":33508}\n\t]}')},R={name:"P1 Preset",tooltip:P,gear:u.fromJsonString('{"items": [\n\t\t{"id":44007,"enchant":3819,"gems":[41401,40017]},\n\t\t{"id":40071},\n\t\t{"id":39719,"enchant":3809,"gems":[39998]},\n\t\t{"id":40723,"enchant":3859},\n\t\t{"id":44002,"enchant":3832,"gems":[39998,40026]},\n\t\t{"id":44008,"enchant":2332,"gems":[39998,0]},\n\t\t{"id":40460,"enchant":3246,"gems":[40017,0]},\n\t\t{"id":40561,"enchant":3601,"gems":[39998]},\n\t\t{"id":40379,"enchant":3719,"gems":[39998,40017]},\n\t\t{"id":40558,"enchant":3606},\n\t\t{"id":40719},\n\t\t{"id":40375},\n\t\t{"id":37111},\n\t\t{"id":40432},\n\t\t{"id":40395,"enchant":3834},\n\t\t{"id":39766},\n\t\t{"id":40342}\n\t]}')},x={name:"P2 Preset",tooltip:P,gear:u.fromJsonString('{"items": [\n\t\t{"id":46184,"enchant":3819,"gems":[41401,39998]},\n\t\t{"id":45243,"gems":[39998]},\n\t\t{"id":46187,"enchant":3809,"gems":[39998]},\n\t\t{"id":45618,"enchant":3831,"gems":[39998]},\n\t\t{"id":45519,"enchant":3832,"gems":[40017,39998,40026]},\n\t\t{"id":45446,"enchant":2332,"gems":[39998,0]},\n\t\t{"id":46183,"enchant":3246,"gems":[39998,0]},\n\t\t{"id":45616,"gems":[39998,39998,39998]},\n\t\t{"id":46185,"enchant":3719,"gems":[40026,39998]},\n\t\t{"id":45135,"enchant":3606,"gems":[39998,40017]},\n\t\t{"id":45495,"gems":[40017]},\n\t\t{"id":45946,"gems":[40017]},\n\t\t{"id":45703},\n\t\t{"id":45535},\n\t\t{"id":46017,"enchant":3834},\n\t\t{"id":45271},\n\t\t{"id":40342}\n\t]}')},A=b({fieldName:"innervateTarget",id:S.fromSpellId(29166),extraCssClasses:["within-raid-sim-hide"],getValue:t=>t.getSpecOptions().innervateTarget?.type==O.Player,setValue:(t,e,n)=>{const a=e.getSpecOptions();a.innervateTarget=r.create({type:n?O.Player:O.Unknown,index:0}),e.setSpecOptions(t,a)}}),H={inputs:[]};class J extends I{constructor(t,e){super(t,e,{cssClass:"restoration-druid-sim-ui",cssScheme:"druid",knownIssues:[],epStats:[y.StatIntellect,y.StatSpirit,y.StatSpellPower,y.StatSpellCrit,y.StatSpellHaste,y.StatMP5],epReferenceStat:y.StatSpellPower,displayStats:[y.StatHealth,y.StatMana,y.StatStamina,y.StatIntellect,y.StatSpirit,y.StatSpellPower,y.StatSpellCrit,y.StatSpellHaste,y.StatMP5],defaults:{gear:R.gear,epWeights:T.fromMap({[y.StatIntellect]:.38,[y.StatSpirit]:.34,[y.StatSpellPower]:1,[y.StatSpellCrit]:.69,[y.StatSpellHaste]:.77,[y.StatMP5]:0}),consumes:W,rotation:E,talents:w.data,specOptions:F,raidBuffs:j,partyBuffs:C,individualBuffs:k,debuffs:B,other:D},playerIconInputs:[A],rotationInputs:H,includeBuffDebuffInputs:[],excludeBuffDebuffInputs:[],otherInputs:{inputs:[v]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[w,G],gear:[M,R,x]}})}}export{w as C,E as D,R as P,J as R,F as a,W as b,x as c};
//# sourceMappingURL=sim-98df7f32.chunk.js.map
